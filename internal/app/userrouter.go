package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

func (app *App) addWebAuthRouter(parent chi.Router) chi.Router {
	return parent.Route("/auth", func(child chi.Router) {
		child.Get("/signup", app.WebEP.InitSignUpUser)
		child.Post("/signup", app.WebEP.SignUpUser)
		child.Get("/signin", app.WebEP.InitSignInUser)
		child.Post("/signin", app.WebEP.SignInUser)
		child.Get("/signout", app.WebEP.SignOutUser)
	})
}

// Thes handlers require authorization
func (app *App) addWebUserRouter(parent chi.Router) chi.Router {
	return parent.Route("/users", func(child chi.Router) {
		child.Use(app.WebEP.ReqAuth)
		child.Get("/", app.WebEP.IndexUsers)
		child.Get("/new", app.WebEP.NewUser)
		child.Post("/", app.WebEP.CreateUser)
		child.Route("/{slug}", func(subChild chi.Router) {
			subChild.Use(userCtx)
			subChild.Get("/", app.WebEP.ShowUser)
			subChild.Get("/edit", app.WebEP.EditUser)
			subChild.Patch("/", app.WebEP.UpdateUser)
			subChild.Put("/", app.WebEP.UpdateUser)
			subChild.Post("/init-delete", app.WebEP.InitDeleteUser)
			subChild.Delete("/", app.WebEP.DeleteUser)
			subChild.Route("/{token}", func(subSubChild chi.Router) {
				subSubChild.Use(confCtx)
				subSubChild.Get("/confirm", app.WebEP.ConfirmUser)
			})
		})
	})
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		ctx := context.WithValue(r.Context(), web.SlugCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func confCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "token")
		ctx := context.WithValue(r.Context(), web.ConfCtxKey, slug)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
