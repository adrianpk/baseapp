package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// Thes handlers require authorization
func (app *App) addWebAccountRouter(parent chi.Router) chi.Router {
	return parent.Route("/accounts", func(uar chi.Router) {
		uar.Use(app.WebEP.ReqAuth)
		//uar.Get("/", app.WebEP.IndexAccounts)
		//uar.Get("/new", app.WebEP.NewAccount)
		//uar.Post("/", app.WebEP.CreateAccount)
		uar.Route("/{slug}", func(uarid chi.Router) {
			uarid.Use(accountCtx)
			//uarid.Get("/", app.WebEP.ShowAccount)
			//uarid.Get("/edit", app.WebEP.EditAccount)
			//uarid.Patch("/", app.WebEP.UpdateAccount)
			//uarid.Put("/", app.WebEP.UpdateAccount)
			//uarid.Post("/init-delete", app.WebEP.InitDeleteAccount)
			//uarid.Delete("/", app.WebEP.DeleteAccount)
			uarid.Get("/roles", app.WebEP.IndexAccountRoles)
		})
	})
}

func accountCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		ctx := context.WithValue(r.Context(), web.SlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}