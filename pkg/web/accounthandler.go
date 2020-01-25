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
	// Info
	RoleAppendedInfoMsg = "role_appended_info_msg"
	RoleRemovedInfoMsg  = "role_removed_info_msg"
)

// IndexAccountRoles web endpoint.
func (ep *Endpoint) IndexAccountRoles(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get account.
	account, err := ep.Service.GetAccount(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all available roles from repo.
	notApplied, err := ep.Service.GetNotAccountRoles(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all account associated roles from repo.
	applied, err := ep.Service.GetAccountRoles(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Group both lists into a map
	l := map[string]interface{}{
		"account":     account.ToForm(),
		"not-applied": model.ToRoleFormList(notApplied),
		"applied":     model.ToRoleFormList(applied),
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

// AppendAccountRole web endpoint.
func (ep *Endpoint) AppendAccountRole(w http.ResponseWriter, r *http.Request) {
	accountSlug, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	accountForm := model.AccountForm{Slug: accountSlug}

	// Decode request data into a form.
	roleForm := model.RoleForm{}
	err = ep.FormToModel(r, &roleForm)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPathRoles(accountForm), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to append role.
	err = ep.Service.AppendAccountRole(accountSlug, roleForm.Slug)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}
	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, RoleAppendedInfoMsg)
	ep.RedirectWithFlash(w, r, AccountPathRoles(accountForm), m, kbs.InfoMT)
}

// RemoveAccountRole web endpoint.
func (ep *Endpoint) RemoveAccountRole(w http.ResponseWriter, r *http.Request) {
	accountSlug, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	accountForm := model.AccountForm{Slug: accountSlug}

	// Decode request data into a form.
	roleForm := model.RoleForm{}
	err = ep.FormToModel(r, &roleForm)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPathRoles(accountForm), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to append role.
	err = ep.Service.RemoveAccountRole(accountSlug, roleForm.Slug)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}
	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, RoleRemovedInfoMsg)
	ep.RedirectWithFlash(w, r, AccountPathRoles(accountForm), m, kbs.InfoMT)
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
