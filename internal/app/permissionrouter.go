package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// These handlers require authorization
func (app *App) addWebPermissionRouter(parent chi.Router) chi.Router {
	return parent.Route("/permissions", func(uar chi.Router) {
		uar.Use(app.WebEP.ReqAuth)
		uar.Get("/", app.WebEP.IndexPermissions)
		uar.Get("/new", app.WebEP.NewPermission)
		uar.Post("/", app.WebEP.CreatePermission)
		uar.Route("/{slug}", func(uarid chi.Router) {
			uarid.Use(permissionCtx)
			uarid.Get("/", app.WebEP.ShowPermission)
			uarid.Get("/edit", app.WebEP.EditPermission)
			uarid.Patch("/", app.WebEP.UpdatePermission)
			uarid.Put("/", app.WebEP.UpdatePermission)
			uarid.Post("/init-delete", app.WebEP.InitDeletePermission)
			uarid.Delete("/", app.WebEP.DeletePermission)
		})
	})
}

func permissionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		ctx := context.WithValue(r.Context(), web.SlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
