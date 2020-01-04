package web

import (
	//"encoding/gob"

	"net/http"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/app/svc"
)

type (
	Endpoint struct {
		*kbs.WebEndpoint
		Service *svc.Service
	}
)

const (
	// Generic
	CannotProcErrMsg  = "cannot_proc_err_msg"
	InputValuesErrMsg = "input_values_err_msg"
	// Fields
	RequiredErrMsg   = "required_err_msg"
	MinLengthErrMsg  = "min_length_err_msg"
	MaxLengthErrMsg  = "max_length_err_msg"
	NotAllowedErrMsg = "not_allowed_err_msg"
	NotEmailErrMsg   = "not_email_err_msg"
	ConfMatchErrMsg  = "conf_match_err_msg"
)

const (
	signedInCookieKey = "signed-in-slug"
)

func NewEndpoint(cfg *kbs.Config, log kbs.Logger, name string) (*Endpoint, error) {
	//registerGobTypes()

	wep, err := kbs.MakeWebEndpoint(cfg, log, pathFxs)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		WebEndpoint: wep,
	}, nil
}

// registerGobTypes
func registerGobTypes() {
	// gob.Register(CustomType1{})
	// gob.Register(CustomType2{})
	// gob.Register(CustomType3{})
}

func (ep *Endpoint) SignedInCookieKey() string {
	return signedInCookieKey
}

// Middlewares
// ReqAuth require user authentication middleware.
func (ep *Endpoint) ReqAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		s, ok := ep.IsAuthenticated(r)
		if !ok {
			ep.Log.Debug("User not authenticated")
			http.Redirect(w, r, AuthPathSignIn(), 302)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		ep.Log.Debug("User authenticated", "slug", s)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (ep *Endpoint) IsAuthenticated(r *http.Request) (slug string, ok bool) {
	return ep.ReadCookieVal(r, ep.SignedInCookieKey())
}
