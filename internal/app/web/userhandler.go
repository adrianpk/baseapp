package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

const (
	userRes = "user"
)

const (
	SlugCtxKey kbs.ContextKey = "slug"
	ConfCtxKey kbs.ContextKey = "conf"
)

const (
	// Defined in 'assets/web/embed/i18n/xx.json'
	UserCreatedInfoID = "user_created_info_msg"
	UserUpdatedInfoID = "user_updated_info_msg"
	UserDeletedInfoID = "user_deleted_info_msg"
	SignedUpInfoID    = "signed_up_info_msg"
	ConfirmedInfoID   = "confirmed_info_msg"
	LoggedInInfoID    = "logged_in_info_msg"
	// Error
	CreateUserErrID  = "create_user_err_msg"
	IndexUsersErrID  = "get_all_users_err_msg"
	GetUserErrID     = "get_user_err_msg"
	UpdateUserErrID  = "update_user_err_msg"
	DeleteUserErrID  = "delete_user_err_msg"
	CredentialsErrID = "credentials_err_msg"
)

// IndexUsers web endpoint.
func (ep *Endpoint) IndexUsers(w http.ResponseWriter, r *http.Request) {
	// Get users list from registered service
	users, err := ep.Service.IndexUsers()
	if err != nil {
		ep.errorRedirect(w, r, "/", IndexUsersErrID, err)
		return
	}

	// Convert result list into a form list
	// Models use sql null types but templates looks
	// clearer if we use plain Go type.
	// i.e.: $user.Username instead of $user.Username.String
	l := model.ToUserFormList(users)
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(userRes, kbs.IndexTmpl)
	if err != nil {
		ep.errorRedirect(w, r, "/", IndexUsersErrID, err)
		return
	}

	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.errorRedirect(w, r, "/", IndexUsersErrID, err)
		return
	}
}

func (ep *Endpoint) NewUser(w http.ResponseWriter, r *http.Request) {
	userForm := model.UserForm{}

	// Wrap response
	wr := ep.WrapRes(w, r, &userForm, nil)
	wr.SetAction(userCreateAction())

	// Get template to render from cache.
	ts, err := ep.TemplateFor(userRes, kbs.NewTmpl)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Write response
	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
		return
	}
}

func (ep *Endpoint) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Decode request data into a form.
	userForm := model.UserForm{}
	err := ep.FormToModel(r, &userForm)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Create a model using form values.
	user := userForm.ToModel()

	// Update non form values
	// NOTE: Use user's IP only on SignUp
	// user.LastIP = db.ToNullString("0.0.0.0/24")

	// Use registered service to do everything related
	// to user creation.
	ves, err := ep.Service.CreateUser(&user)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.rerenderUserForm(w, r, user.ToForm(), ves, kbs.NewTmpl, userCreateAction())
		return
	}

	// Then take care of other kind of possible errors
	// that service can generate.
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), CreateUserErrID, err)
		return
	}

	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.localize(r, UserCreatedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, kbs.InfoMT)
}

// ShowUser web endpoint.
func (ep *Endpoint) ShowUser(w http.ResponseWriter, r *http.Request) {
	// Get slug from request context.
	s, err := ep.getSlug(r)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Use registered service to do everything related
	// to user creation.
	user, err := ep.Service.GetUser(s)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Wrap response
	wr := ep.WrapRes(w, r, user.ToForm(), nil)

	// Template
	ts, err := ep.TemplateFor(userRes, kbs.ShowTmpl)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}
}

// EditUser web endpoint.
func (ep *Endpoint) EditUser(w http.ResponseWriter, r *http.Request) {
	// Identifier
	s, err := ep.getSlug(r)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Use registerd service to get the user from repo.
	user, err := ep.Service.GetUser(s)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Wrap response
	userForm := user.ToForm()
	wr := ep.WrapRes(w, r, &userForm, nil)
	wr.SetAction(userUpdateAction(&userForm))

	// Template
	ts, err := ep.TemplateFor(userRes, kbs.EditTmpl)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}
}

// UpdateUser web endpoint.
func (ep *Endpoint) UpdateUser(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Decode request data into a form.
	userForm := model.UserForm{}
	err = ep.FormToModel(r, &userForm)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	// Create a model using form values.
	user := userForm.ToModel()

	// Use registered service to do everything related
	// to user update.
	ves, err := ep.Service.UpdateUser(s, &user)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.Log.Warn("Validation errors", "dump", fmt.Sprintf("%+v", ves.FieldErrors))
		ep.rerenderUserForm(w, r, user.ToForm(), ves, kbs.NewTmpl, userCreateAction())
		return
	}

	// Non validation errors
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), UpdateUserErrID, err)
		return
	}

	m := ep.localize(r, UserUpdatedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, kbs.InfoMT)
}

// InitDeleteUser web endpoint.
func (ep *Endpoint) InitDeleteUser(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Use registerd service to get the user from repo.
	user, err := ep.Service.GetUser(s)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), GetUserErrID, err)
		return
	}

	// Wrap response
	userForm := user.ToForm()
	wr := ep.WrapRes(w, r, &userForm, nil)
	wr.SetAction(userDeleteAction(&userForm))

	// Template
	ts, err := ep.TemplateFor(userRes, kbs.InitDelTmpl)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), DeleteUserErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), DeleteUserErrID, err)
		return
	}
}

// DeleteUser web endpoint.
func (ep *Endpoint) DeleteUser(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), DeleteUserErrID, err)
		return
	}

	// Service
	err = ep.Service.DeleteUser(s)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), DeleteUserErrID, err)
		return
	}

	m := ep.localize(r, UserDeletedInfoID)
	ep.RedirectWithFlash(w, r, UserPath(), m, kbs.InfoMT)
}

//func (ep *Endpoint) InitSignUpUser(w http.ResponseWriter, r *http.Request) {
//// Req & Res
//res := &tp.SignUpUserRes{}
//res.Action = s.userSignUpAction()

//// Wrap response
//wr := s.OKRes(w, r, res, "")

//// Template
//ts, err := s.TemplateFor(userRes, kbs.SignUpTmpl)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
//return
//}

//// Write response
//err = ts.Execute(w, wr)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
//return
//}
//}

//// SignUpUser web endpoint.
//func (ep *Endpoint) SignUpUser(w http.ResponseWriter, r *http.Request) {
//var req tp.SignUpUserReq
//var res tp.SignUpUserRes
//res.IsNew = true
//res.Action = s.userSignUpAction()

//// Input data to request struct
//err := s.FormToModel(r, &req.User)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
//return
//}

//// Service
//err = s.service.SignUpUser(req, &res)

//// Input validation errors
//if !res.Errors.IsEmpty() {
//s.rerenderUserForm(w, r, res, kbs.SignUpTmpl)
//return
//}

//// Non validation errors
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CreateUserErrID, err)
//return
//}

//m := s.localize(r, SignedUpInfoID)
//s.RedirectWithFlash(w, r, UserPath(), m, kbs.InfoMT)
//}

//func (ep *Endpoint) InitSignInUser(w http.ResponseWriter, r *http.Request) {
//// Req & Res
//res := &tp.SignInUserRes{}
//res.Action = s.userSignInAction()

//// Wrap response
//wr := s.OKRes(w, r, res, "")

//// Template
//ts, err := s.TemplateFor(userRes, kbs.SignInTmpl)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
//return
//}

//// Write response
//err = ts.Execute(w, wr)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
//return
//}
//}

//// ConfirmUser web endpoint.
//func (ep *Endpoint) ConfirmUser(w http.ResponseWriter, r *http.Request) {
//var req tp.GetUserReq
//var res tp.GetUserRes

//// Identifier
//slug, err := s.getUserSlug(r)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), GetUserErrID, err)
//return
//}

//// Token
//token, err := s.getToken(r)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), GetUserErrID, err)
//return
//}

//req = tp.GetUserReq{
//tp.Identifier{
//Slug:  slug,
//Token: token,
//},
//}

//// Service
//err = s.service.ConfirmUser(req, &res)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), GetUserErrID, err)
//return
//}

//m := s.localize(r, UserCreatedInfoID)
//s.RedirectWithFlash(w, r, UserPath(), m, kbs.InfoMT)
//}

//// SignInUser web endpoint.
//func (ep *Endpoint) SignInUser(w http.ResponseWriter, r *http.Request) {
//var req tp.SignInUserReq
//var res tp.SignInUserRes

//// Input data to request struct
//err := s.FormToModel(r, &req.SignIn)
//if err != nil {
//s.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
//return
//}

//// Service
//err = s.service.SignInUser(req, &res)
//if err != nil {
//s.errorRedirect(w, r, UserPathSignIn(), CredentialsErrID, err)
//return
//}

//m := s.localize(r, LoggedInInfoID)
//s.RedirectWithFlash(w, r, UserPath(), m, kbs.InfoMT)
//}

func (ep *Endpoint) rerenderUserForm(w http.ResponseWriter, r *http.Request, data interface{}, valErrors kbs.ValErrorSet, template string, action kbs.FormAction) {
	wr := ep.WrapRes(w, r, data, valErrors)
	wr.AddErrorFlash(InputValuesErrID)
	wr.SetAction(action)

	ts, err := ep.TemplateFor(userRes, template)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), InputValuesErrID, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.errorRedirect(w, r, UserPath(), CannotProcErrID, err)
		return
	}

	return
}

// Localization - I18N
func (ep *Endpoint) localize(r *http.Request, msgID string) string {
	l := ep.Localizer(r)
	if l == nil {
		ep.Log.Warn("No localizer available")
		return msgID
	}

	t, _, err := l.LocalizeWithTag(&i18n.LocalizeConfig{
		MessageID: msgID,
	})

	if err != nil {
		ep.Log.Error(err)
		return msgID
	}

	//s.Log.Debug("Localized message", "value", t, "lang", lang)

	return t
}

func (ep *Endpoint) localizeMessageID(l *i18n.Localizer, messageID string) (string, error) {
	return l.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// Misc
func (ep *Endpoint) getSlug(r *http.Request) (slug string, err error) {
	ctx := r.Context()
	slug, ok := ctx.Value(SlugCtxKey).(string)
	if !ok {
		err := errors.New("no slug provided")
		return "", err
	}

	return slug, nil
}

func (ep *Endpoint) getToken(r *http.Request) (token string, err error) {
	ctx := r.Context()
	token, ok := ctx.Value(ConfCtxKey).(string)
	if !ok {
		err := errors.New("no token provided")
		return "", err
	}

	return token, nil
}

// userCreateAction
func userCreateAction() kbs.FormAction {
	return kbs.FormAction{Target: fmt.Sprintf("%s", UserPath()), Method: "POST"}
}

// userUpdateAction
func userUpdateAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: UserPathSlug(model), Method: "PUT"}
}

// userDeleteAction
func userDeleteAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: UserPathSlug(model), Method: "DELETE"}
}

// userSignUpAction
func userSignUpAction() kbs.FormAction {
	return kbs.FormAction{Target: UserPathSignUp(), Method: "POST"}
}

// userSignInAction
func userSignInAction() kbs.FormAction {
	return kbs.FormAction{Target: UserPathSignIn(), Method: "POST"}
}

func (ep *Endpoint) errorRedirect(w http.ResponseWriter, r *http.Request, redirPath, msgID string, err error) {
	m := ep.localize(r, msgID)
	ep.RedirectWithFlash(w, r, redirPath, m, kbs.ErrorMT)
	ep.Log.Error(err)
}
