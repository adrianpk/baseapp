package web

import (
	//"encoding/gob"

	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/svc"
)

type (
	Endpoint struct {
		*kbs.WebEndpoint
		service   *svc.Service
		authCache AuthCache
	}
)

const (
	SlugCtxKey    kbs.ContextKey = "slug"
	SubSlugCtxKey kbs.ContextKey = "subslug"
	ConfCtxKey    kbs.ContextKey = "conf"
)

const (
	// Info
	CreatedInfoMsg = "created_info_msg"
	UpdatedInfoMsg = "updated_info_msg"
	DeletedInfoMsg = "deleted_info_msg"
	// Error
	CreateErrMsg = "create_err_msg"
	IndexErrMsg  = "get_all_err_msg"
	GetErrMsg    = "get_err_msg"
	UpdateErrMsg = "update_err_msg"
	DeleteErrMsg = "delete_err_msg"
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

func NewEndpoint(cfg *kbs.Config, log kbs.Logger, name string) (*Endpoint, error) {
	//registerGobTypes()
	wep, err := kbs.MakeWebEndpoint(cfg, log, pathFxs)
	if err != nil {
		return nil, err
	}

	return &Endpoint{
		WebEndpoint: wep,
		authCache:   NewCache(),
	}, nil
}

// registerGobTypes
func registerGobTypes() {
	// gob.Register(CustomType1{})
	// gob.Register(CustomType2{})
	// gob.Register(CustomType3{})
}

func (ep *Endpoint) Service() *svc.Service {
	return ep.service
}

func (ep *Endpoint) SetService(s *svc.Service) {
	ep.service = s
	ep.authCache.SetService(s)
}

// Middlewares
// ReqAuth require user authentication middleware.
func (ep *Endpoint) ReqAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ad, ok := ep.IsAuthenticated(r)

		// Not authenticated
		if !ok {
			ep.Log.Debug("User not authenticated")
			http.Redirect(w, r, AuthPathSignIn(), 302)
			return
		}

		// Authenticated
		username, ok0 := ad["username"]
		slug, ok1 := ad["slug"]

		// If superadmin don't ask for other permissions
		if ok0 && ok1 &&
			username == "superadmin" &&
			slug == "superadmin-000000000002" {

			w.Header().Add("Cache-Control", "no-store")

			ep.Log.Debug("User superadmin authenticated")

			next.ServeHTTP(w, r)
			return
		}

		userTags, ok := ad["permissions"]

		// Cannot get associated permissions
		if !ok {
			ep.Log.Debug("Cannot get user account permissions")

			http.Redirect(w, r, AuthPathSignIn(), 302)
			return
		}

		// Get user permissions
		ep.Log.Debug("Account permissionis", "tags", spew.Sdump(userTags))

		// Get request path (resource)
		path := filepath.Clean(r.URL.Path)
		ep.Log.Debug("Request", "path", path)

		// Get required permissions to access this path (resource)
		reqTags, err := ep.authCache.PathPermissionTags(path)
		if err != nil {
			ep.Log.Error(err, "Cannot get required resource permission tags")
			ep.Log.Warn("Account not authorized", "username", username, "path", path)
			http.Redirect(w, r, AuthPathSignIn(), 302)
			return
		}

		ep.Log.Debug("Resource", "required-permission-tags", spew.Sdump(reqTags))

		// Verify that user has all the required permissions
		for _, t := range reqTags {
			if !strings.Contains(userTags, t) {
				ep.Log.Error(err, "Cannot get required resource permission tags")
				ep.Log.Info("Account not authorized", "required", t)
				http.Redirect(w, r, AuthPathSignIn(), 302)
				return
			}
		}

		// User is authorized
		w.Header().Add("Cache-Control", "no-store")

		ep.Log.Debug("User authenticated", "slug", ad)

		next.ServeHTTP(w, r)

	}

	return http.HandlerFunc(fn)
}

func (ep *Endpoint) getSlug(r *http.Request) (slug string, err error) {
	ctx := r.Context()
	slug, ok := ctx.Value(SlugCtxKey).(string)
	if !ok {
		err := errors.New("no slug provided")
		return "", err
	}

	return slug, nil
}

func (ep *Endpoint) getSubSlug(r *http.Request) (slug string, err error) {
	ctx := r.Context()
	slug, ok := ctx.Value(SubSlugCtxKey).(string)
	if !ok {
		err := errors.New("no subslug provided")
		return "", err
	}

	return slug, nil
}

func (ep *Endpoint) ErrorRedirect(w http.ResponseWriter, r *http.Request, redirPath, messageID string, err error) {
	m := ep.Localize(r, messageID)
	ep.RedirectWithFlash(w, r, redirPath, m, kbs.ErrorMT)
	ep.Log.Error(err)
}
