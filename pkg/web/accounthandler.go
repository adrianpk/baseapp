package web

import (
	"fmt"
	"net/http"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

const (
	accountRes = "account"
)

const (
	accountResPl = "accounts"
)

const (
	RolesTmpl = "roles.tmpl"
)

// IndexAccountRoles web endpoint.
func (ep *Endpoint) IndexAccountRoles(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get all available roles from repo.
	all, err := ep.Service.IndexRoles()
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all account associated roles from repo.
	accountRoles, err := ep.Service.GetAccountRoles(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Group both lists into a map
	l := map[string][]model.RoleForm{
		"all":     model.ToRoleFormList(all),
		"account": model.ToRoleFormList(accountRoles),
	}

	// Wrap response
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(accountRes, RolesTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, "/", CannotProcErrMsg, err)
		return
	}

	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, "/", CannotProcErrMsg, err)
		return
	}
}

func (ep *Endpoint) rerenderAccountForm(w http.ResponseWriter, r *http.Request, data interface{}, valErrors kbs.ValErrorSet, template string, action kbs.FormAction) {
	wr := ep.WrapRes(w, r, data, valErrors)
	wr.AddErrorFlash(InputValuesErrMsg)
	wr.SetAction(action)

	ts, err := ep.TemplateFor(accountRes, template)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), InputValuesErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	return
}

// accountCreateAction
func accountCreateAction() kbs.FormAction {
	return kbs.FormAction{Target: fmt.Sprintf("%s", AccountPath()), Method: "POST"}
}

// accountUpdateAction
func accountUpdateAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: AccountPathSlug(model), Method: "PUT"}
}

// accountDeleteAction
func accountDeleteAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: AccountPathSlug(model), Method: "DELETE"}
}

// accountSignUpAction
func accountSignUpAction() kbs.FormAction {
	return kbs.FormAction{Target: AuthPathSignUp(), Method: "POST"}
}

// accountSignInAction
func accountSignInAction() kbs.FormAction {
	return kbs.FormAction{Target: AuthPathSignIn(), Method: "POST"}
}

func (ep *Endpoint) eErrorRedirect(w http.ResponseWriter, r *http.Request, redirPath, msgID string, err error) {
	m := ep.Localize(r, msgID)
	ep.RedirectWithFlash(w, r, redirPath, m, kbs.ErrorMT)
	ep.Log.Error(err)
}
