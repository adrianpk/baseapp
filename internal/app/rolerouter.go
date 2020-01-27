package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// These handlers require authorization
func (app *App) addWebRoleRouter(parent chi.Router) chi.Router {
	return parent.Route("/roles", func(child chi.Router) {
		child.Use(app.WebEP.ReqAuth)
		child.Get("/", app.WebEP.IndexRoles)
		child.Get("/new", app.WebEP.NewRole)
		child.Post("/", app.WebEP.CreateRole)
		child.Route("/{slug}", func(subChild chi.Router) {
			subChild.Use(roleCtx)
			subChild.Get("/", app.WebEP.ShowRole)
			subChild.Get("/edit", app.WebEP.EditRole)
			subChild.Patch("/", app.WebEP.UpdateRole)
			subChild.Put("/", app.WebEP.UpdateRole)
			subChild.Post("/init-delete", app.WebEP.InitDeleteRole)
			subChild.Delete("/", app.WebEP.DeleteRole)
			subChild.Get("/permissions", app.WebEP.IndexRolePermissions)
			subChild.Post("/permissions", app.WebEP.AppendRolePermission)
			subChild.Route("/permissions/{subSlug}", func(subSubChild chi.Router) {
				subSubChild.Use(rolePermissionCtx)
				subSubChild.Delete("/", app.WebEP.RemoveRolePermission)
			})
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

func rolePermissionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "subSlug")
		ctx := context.WithValue(r.Context(), web.SubSlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
