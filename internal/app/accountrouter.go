package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// Thes handlers require authorization
func (app *App) addWebAccountRouter(parent chi.Router) chi.Router {
	return parent.Route("/accounts", func(child chi.Router) {
		child.Use(app.WebEP.ReqAuth)
		//child.Get("/", app.WebEP.IndexAccounts)
		//child.Get("/new", app.WebEP.NewAccount)
		//child.Post("/", app.WebEP.CreateAccount)
		child.Route("/{slug}", func(subChild chi.Router) {
			subChild.Use(accountCtx)
			//subChild.Get("/", app.WebEP.ShowAccount)
			//subChild.Get("/edit", app.WebEP.EditAccount)
			//subChild.Patch("/", app.WebEP.UpdateAccount)
			//subChild.Put("/", app.WebEP.UpdateAccount)
			//subChild.Post("/init-delete", app.WebEP.InitDeleteAccount)
			//subChild.Delete("/", app.WebEP.DeleteAccount)
			subChild.Get("/accountroles", app.WebEP.IndexAccountRoles)
			subChild.Post("/accountroles", app.WebEP.AppendAccountRole)
			subChild.Route("/accountroles/{subSlug}", func(subSubChild chi.Router) {
				subSubChild.Use(accountRoleCtx)
				subSubChild.Delete("/", app.WebEP.RemoveAccountRole)
			})
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

func accountRoleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "subSlug")
		ctx := context.WithValue(r.Context(), web.SubSlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
