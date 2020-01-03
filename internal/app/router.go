package app

import (
	"net/http"
	"os"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"github.com/markbates/pkger"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type (
	textResponse string
)

var (
	langMatcher = language.NewMatcher(message.DefaultCatalog.Languages())
)

func (app *App) NewWebRouter() *kbs.Router {
	rt := app.makeWebHomeRouter(app.Cfg, app.Log)
	app.addWebUserRouter(rt)
	return rt
}

func (app *App) makeWebHomeRouter(cfg *kbs.Config, log kbs.Logger) *kbs.Router {
	rt := kbs.NewRouter(cfg, log, "web-home-router")
	app.addWebHomeRoutes(rt)
	return rt
}

func (app *App) addWebHomeRoutes(rt *kbs.Router) {
	dir := "/assets/web/embed/public"
	fs := http.FileServer(kbs.FileSystem{pkger.Dir(dir)})

	rt.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := pkger.Stat(dir + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)

		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func (t textResponse) write(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(t))
}
