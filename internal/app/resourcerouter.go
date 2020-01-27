package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// These handlers require authorization
func (app *App) addWebResourceRouter(parent chi.Router) chi.Router {
	return parent.Route("/resources", func(child chi.Router) {
		child.Use(app.WebEP.ReqAuth)
		child.Get("/", app.WebEP.IndexResources)
		child.Get("/new", app.WebEP.NewResource)
		child.Post("/", app.WebEP.CreateResource)
		child.Route("/{slug}", func(subChild chi.Router) {
			subChild.Use(resourceCtx)
			subChild.Get("/", app.WebEP.ShowResource)
			subChild.Get("/edit", app.WebEP.EditResource)
			subChild.Patch("/", app.WebEP.UpdateResource)
			subChild.Put("/", app.WebEP.UpdateResource)
			subChild.Post("/init-delete", app.WebEP.InitDeleteResource)
			subChild.Delete("/", app.WebEP.DeleteResource)
			subChild.Get("/permissions", app.WebEP.IndexResourcePermissions)
			subChild.Post("/permissions", app.WebEP.AppendResourcePermission)
			subChild.Route("/permissions/{subSlug}", func(subSubChild chi.Router) {
				subSubChild.Use(resourcePermissionCtx)
				subSubChild.Delete("/", app.WebEP.RemoveResourcePermission)
			})
		})
	})
}

func resourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		ctx := context.WithValue(r.Context(), web.SlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func resourcePermissionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "subSlug")
		ctx := context.WithValue(r.Context(), web.SubSlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
