package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// These handlers require authorization
func (app *App) addWebRoleRouter(parent chi.Router) chi.Router {
	return parent.Route("/roles", func(uar chi.Router) {
		uar.Use(app.WebEP.ReqAuth)
		uar.Get("/", app.WebEP.IndexRoles)
		uar.Get("/new", app.WebEP.NewRole)
		uar.Post("/", app.WebEP.CreateRole)
		uar.Route("/{slug}", func(uarid chi.Router) {
			uarid.Use(roleCtx)
			uarid.Get("/", app.WebEP.ShowRole)
			uarid.Get("/edit", app.WebEP.EditRole)
			uarid.Patch("/", app.WebEP.UpdateRole)
			uarid.Put("/", app.WebEP.UpdateRole)
			uarid.Post("/init-delete", app.WebEP.InitDeleteRole)
			uarid.Delete("/", app.WebEP.DeleteRole)
		})
	})
}

func roleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		ctx := context.WithValue(r.Context(), web.SlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
