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
		Service *svc.Service
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
	}, nil
}

// registerGobTypes
func registerGobTypes() {
	// gob.Register(CustomType1{})
	// gob.Register(CustomType2{})
	// gob.Register(CustomType3{})
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

		// NOTE: WIP, This is a non-optimized implementation
		// of logic that allows you to allow/restrict access
		// resourcers. A more stright access to persistence
		// can be implemented if required.
		// A cache layer can also be helpful.
		username, ok0 := s["username"]
		slug, ok1 := s["slug"]

		// If superadmin don't ask for other permissions
		if ok0 && ok1 &&
			username == "superadmin" &&
			slug == "superadmin-000000000002" {

			w.Header().Add("Cache-Control", "no-store")

			ep.Log.Debug("User superadm authenticated")

			next.ServeHTTP(w, r)
			return
		}

		userTags, ok := s["permissions"]
		if !ok {
			ep.Log.Debug("User not authenticated")
			http.Redirect(w, r, AuthPathSignIn(), 302)
			return
		}

		// Get user permissions
		ep.Log.Info("User permission", "tags", spew.Sdump(userTags))

		// Get request path (resource)
		path := filepath.Clean(r.URL.Path)
		ep.Log.Info("Request", "path", path)

		// Get required permissions to access this path (resource)
		reqTags, err := ep.Service.ResourcePermissionTagsByPath(path)
		if err != nil {
			ep.Log.Debug("User not authorized")
			ep.Log.Error(err, "Cannot get required resource permission tags")
			http.Redirect(w, r, AuthPathSignIn(), 302)
			return
		}

		ep.Log.Info("Resource", "required-permission-tags", spew.Sdump(reqTags))

		// Verify that user has all the required permissions
		for _, t := range reqTags {
			if !strings.Contains(userTags, t) {
				ep.Log.Info("User not authorized", "required-not-found", t)
				ep.Log.Error(err, "Cannot get required resource permission tags")
				http.Redirect(w, r, AuthPathSignIn(), 302)
				return
			}
		}

		// User is authorized
		w.Header().Add("Cache-Control", "no-store")

		ep.Log.Debug("User authenticated", "slug", s)

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
