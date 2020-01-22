package app

import (
	"fmt"
	"net/http"
	"sync"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

type (
	App struct {
		*kbs.App
		WebEP *web.Endpoint
	}
)

// NewApp creates a new app  worker instance.func NewApp(cfg *kbs.Config, log kbs.Logger, name string, core *kbs.Worker) (*App, error) {
func NewApp(cfg *kbs.Config, log kbs.Logger, name string) (*App, error) {
	app := App{
		App: kbs.NewApp(cfg, log, name),
	}

	// Endpoint
	wep, err := web.NewEndpoint(cfg, log, "web-endpoint")
	if err != nil {
		return nil, err
	}
	app.WebEP = wep

	// Router
	app.WebRouter = app.NewWebRouter()

	return &app, nil
}

// Init runs pre Start process.
func (app *App) Init() error {
	err := app.Migrator.Migrate()
	if err != nil {
		return err
	}

	err = app.Seeder.Seed()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) Start() error {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		app.StartWeb()

		wg.Done()
	}()

	//wg.Add(1)
	//go func() {
	//a.StartJSONREST()
	//wg.Done()
	//}()

	wg.Wait()
	return nil
}

func (app *App) Stop() {
	// TODO: Gracefully stop all workers
}

func (app *App) StartWeb() error {
	p := app.Cfg.ValOrDef("web.server.port", "8080")
	p = fmt.Sprintf(":%s", p)

	app.Log.Info("Web server initializing", "port", p)

	err := http.ListenAndServe(p, app.WebRouter)
	app.Log.Error(err)

	return err
}
