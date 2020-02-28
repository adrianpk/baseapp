package web_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	kbs "gitlab.com/kabestan/backend/kabestan"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/app"
	"gitlab.com/kabestan/repo/baseapp/internal/mig"
	repo "gitlab.com/kabestan/repo/baseapp/internal/repo/pg"
	"gitlab.com/kabestan/repo/baseapp/internal/seed"
	"gitlab.com/kabestan/repo/baseapp/internal/svc"
)

type (
	testSetup struct {
		Cfg      *kbs.Config
		Log      kbs.Logger
		App      *app.App
		DB       *sqlx.DB
		UserRepo repo.UserRepo
		Svc      *svc.Service
	}
)

// RequestMethod - Valid request methods
type RequestMethod string

const (
	// GET request method
	GET RequestMethod = "GET"
	// POST request method
	POST RequestMethod = "POST"
	// PUT request method
	PUT RequestMethod = "PUT"
	// PATCH request method
	PATCH RequestMethod = "PATCH"
	// DELETE request method
	DELETE RequestMethod = "DELETE"
)

var (
	ts         testSetup
	currentEnv = "test"
	ctHeader   = "Content-Type"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

// TestMock does not test anything.
func TestMock(t *testing.T) {
	// Only used to test migration and rollback operations
	// on setup and teardown respectivelly.
}

func setContentType(req *http.Request, contentType string) {
	req.Header.Set(ctHeader, contentType)
}

func buildURL(path string, segments ...string) string {
	u := fmt.Sprintf("http://%s/%s", ts.GetWebServerAddress(), path)
	for _, s := range segments {
		u = fmt.Sprintf("%s/%s", u, s)
	}
	return u
}

func buildResURL(path string, id uuid.UUID) string {
	return fmt.Sprintf("%s/%s", buildURL(path), id)
}

func buildResEditURL(path string, id uuid.UUID) string {
	return fmt.Sprintf("%s/%s/edit", buildURL(path), id)
}

func buildResInitDeleteURL(path string, id uuid.UUID) string {
	return fmt.Sprintf("%s/%s/init-delete", buildURL(path), id)
}

func post(req *http.Request) (*http.Response, error) {
	hc := http.Client{}
	return hc.Do(req)
}

func responseBody(res *http.Response) string {
	resBuf := new(bytes.Buffer)
	resBuf.ReadFrom(res.Body)
	defer res.Body.Close()
	return resBuf.String()
}

func recoverControl() {
	if r := recover(); r != nil {
		ts.Log.Info("Recovered.")
	}
}

func matchString(toMatch, body string) bool {
	//fmt.Println("To match:", toMatch)
	//fmt.Println("Body", body)
	rs := fmt.Sprintf("(%s)", toMatch)
	r, _ := regexp.Compile(rs)
	return r.MatchString(body)
}

func buildRequest(targetURL string, rm RequestMethod) (req *http.Request, err error) {
	req = httptest.NewRequest(string(rm), targetURL, nil)
	req.RequestURI = ""
	req.URL, err = url.Parse(targetURL)
	if err != nil {
		return req, err
	}
	return req, nil
}

func makeFormPostReq(targetURL string, formData url.Values) (req *http.Request, err error) {
	formSubmitCT := "application/x-www-form-urlencoded"
	encodedReader := strings.NewReader(formData.Encode())
	req = httptest.NewRequest("POST", targetURL, encodedReader)
	req.RequestURI = ""
	req.URL, err = url.Parse(targetURL)
	if err != nil {
		return req, err
	}
	setContentType(req, formSubmitCT)
	return req, nil
}

func makeFormPutReq(targetURL string, formData url.Values) (req *http.Request, err error) {
	formSubmitCT := "application/x-www-form-urlencoded"
	encodedReader := strings.NewReader(formData.Encode())
	req = httptest.NewRequest("PUT", targetURL, encodedReader)
	req.RequestURI = ""
	req.URL, err = url.Parse(targetURL)
	if err != nil {
		return req, err
	}
	setContentType(req, formSubmitCT)
	return req, nil
}

func executeRequest(r *http.Request) (res *http.Response, err error) {
	//Configure TLS, etc.
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{
	// 		InsecureSkipVerify: true,
	// 	},
	// }
	hc := http.Client{
		// Transport: tr,
		Timeout: 5 * time.Second,
	}
	// r.Header.Add("User-Agent", "Golang Test")
	res, err = hc.Do(r)
	if err != nil {
		return res, err
	}
	return res, nil
}

func extractBody(r *http.Response) (body string, err error) {
	resBuf := new(bytes.Buffer)
	_, err = resBuf.ReadFrom(r.Body)
	defer r.Body.Close()
	if err != nil {
		return body, err
	}
	return resBuf.String(), nil
}

// Setupt & Teardown
func setup() {
	ts, err := newTestSetup()
	if err != nil {

	}
	// ts.App.Migrator.Reset()
	// ts.App.Migrator.RollbackAll()
	ts.App.Migrator.Migrate()
}

func teardown() {
	ts.App.Migrator.RollbackAll()
}

func newTestSetup() (ts *testSetup, err error) {
	cfg := testConfig()
	log := kbs.NewLogger(cfg)

	a, err := app.NewApp(cfg, log, "app")
	if err != nil {
		return ts, err
	}

	db, err := kbs.NewPostgresConn(cfg, log)
	if err != nil {
		return ts, err
	}

	migrator, err := mig.NewMigrator(cfg, log, "migrator", db)
	if err != nil {
		return ts, err
	}

	seeder, err := seed.NewSeeder(cfg, log, "seeder", db)
	if err != nil {
		return ts, err
	}

	userRepo := repo.NewUserRepo(cfg, log, "user-repo", db)

	svc := svc.NewService(cfg, log, "core-service", db)

	svc.UserRepo = userRepo

	a.Migrator = migrator
	a.Seeder = seeder
	a.WebEP.SetService(svc)
	a.Init()
	a.Start()

	return &testSetup{
		Cfg: cfg,
		Log: log,
		App: a,
		DB:  db,
		Svc: svc,
	}, nil
}

func testConfig() *kbs.Config {
	cfg := &kbs.Config{}
	values := map[string]string{
		"pg.host":               "postgres",
		"pg.port":               "5432",
		"pg.schema":             "public",
		"pg.database":           "baseapp_test",
		"pg.user":               "baseapp",
		"pg.password":           "baseapp",
		"pg.backoff.maxentries": "3",
	}

	cfg.SetNamespace("kbs")
	cfg.SetValues(values)
	return cfg
}

// Type methods
func (ts *testSetup) GetWebServerAddress() string {
	p := ts.Cfg.ValOrDef("web.server.port", "")
	return fmt.Sprintf("%s:%s", "localhost", p)
}
