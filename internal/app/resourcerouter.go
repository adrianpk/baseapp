package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

// These handlers require authorization
func (app *App) addWebResourceRouter(parent chi.Router) chi.Router {
	return parent.Route("/resources", func(uar chi.Router) {
		uar.Use(app.WebEP.ReqAuth)
		uar.Get("/", app.WebEP.IndexResources)
		uar.Get("/new", app.WebEP.NewResource)
		uar.Post("/", app.WebEP.CreateResource)
		uar.Route("/{slug}", func(uarid chi.Router) {
			uarid.Use(resourceCtx)
			uarid.Get("/", app.WebEP.ShowResource)
			uarid.Get("/edit", app.WebEP.EditResource)
			uarid.Patch("/", app.WebEP.UpdateResource)
			uarid.Put("/", app.WebEP.UpdateResource)
			uarid.Post("/init-delete", app.WebEP.InitDeleteResource)
			uarid.Delete("/", app.WebEP.DeleteResource)
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
