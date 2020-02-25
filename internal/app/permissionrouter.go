package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// These handlers require authorization
func (app *App) addWebPermissionRouter(parent chi.Router) chi.Router {
	return parent.Route("/permissions", func(child chi.Router) {
		child.Use(app.WebEP.ReqAuth)
		child.Get("/", app.WebEP.IndexPermissions)
		child.Get("/new", app.WebEP.NewPermission)
		child.Post("/", app.WebEP.CreatePermission)
		child.Route("/{slug}", func(subChild chi.Router) {
			subChild.Use(permissionCtx)
			subChild.Get("/", app.WebEP.ShowPermission)
			subChild.Get("/edit", app.WebEP.EditPermission)
			subChild.Patch("/", app.WebEP.UpdatePermission)
			subChild.Put("/", app.WebEP.UpdatePermission)
			subChild.Post("/init-delete", app.WebEP.InitDeletePermission)
			subChild.Delete("/", app.WebEP.DeletePermission)
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
