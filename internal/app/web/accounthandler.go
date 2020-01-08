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
	// Defined in 'assets/web/embed/i18n/xx.json'
	AccountCreatedInfoMsg = "account_created_info_msg"
	AccountUpdatedInfoMsg = "account_updated_info_msg"
	AccountDeletedInfoMsg = "account_deleted_info_msg"
	// Error
	CreateAccountErrMsg = "create_account_err_msg"
	IndexAccountsErrMsg = "get_all_accounts_err_msg"
	GetAccountErrMsg    = "get_account_err_msg"
	GetAccountsErrMsg   = "get_accounts_err_msg"
	UpdateAccountErrMsg = "update_account_err_msg"
	DeleteAccountErrMsg = "delete_account_err_msg"
)

// IndexAccounts web endpoint.
func (ep *Endpoint) IndexAccounts(w http.ResponseWriter, r *http.Request) {
	// Get accounts list from registered service
	accounts, err := ep.Service.IndexAccounts()
	if err != nil {
		ep.ErrorRedirect(w, r, "/", IndexAccountsErrMsg, err)
		return
	}

	// Convert result list into a form list
	// Models use sql null types but templates looks
	// clearer if we use plain Go type.
	// i.e.: $account.Accountname instead of $account.Accountname.String
	l := model.ToAccountFormList(accounts)
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(accountRes, kbs.IndexTmpl)
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

func (ep *Endpoint) NewAccount(w http.ResponseWriter, r *http.Request) {
	accountForm := model.AccountForm{}

	// Wrap response
	wr := ep.WrapRes(w, r, &accountForm, nil)
	wr.SetAction(accountCreateAction())

	// Get template to render from cache.
	ts, err := ep.TemplateFor(accountRes, kbs.NewTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}
}

func (ep *Endpoint) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// Decode request data into a form.
	accountForm := model.AccountForm{}
	err := ep.FormToModel(r, &accountForm)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	account := accountForm.ToModel()

	// Use registered service to do everything related
	// to account creation.
	ves, err := ep.Service.CreateAccount(&account)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.rerenderAccountForm(w, r, account.ToForm(), ves, kbs.NewTmpl, accountCreateAction())
		return
	}

	// Then take care of other kind of possible errors
	// that service can generate.
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, AccountCreatedInfoMsg)
	ep.RedirectWithFlash(w, r, AccountPath(), m, kbs.InfoMT)
}

// ShowAccount web endpoint.
func (ep *Endpoint) ShowAccount(w http.ResponseWriter, r *http.Request) {
	// Get slug from request context.
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Use registered service to do everything related
	// to account creation.
	account, err := ep.Service.GetAccount(s)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), GetAccountErrMsg, err)
		return
	}

	// Wrap response
	wr := ep.WrapRes(w, r, account.ToForm(), nil)

	// Template
	ts, err := ep.TemplateFor(accountRes, kbs.ShowTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}
}

// EditAccount web endpoint.
func (ep *Endpoint) EditAccount(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the account from repo.
	account, err := ep.Service.GetAccount(s)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), GetAccountErrMsg, err)
		return
	}

	// Wrap response
	accountForm := account.ToForm()
	wr := ep.WrapRes(w, r, &accountForm, nil)
	wr.SetAction(accountUpdateAction(&accountForm))

	// Template
	ts, err := ep.TemplateFor(accountRes, kbs.EditTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}
}

// UpdateAccount web endpoint.
func (ep *Endpoint) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), GetAccountErrMsg, err)
		return
	}

	// Decode request data into a form.
	accountForm := model.AccountForm{}
	err = ep.FormToModel(r, &accountForm)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	account := accountForm.ToModel()

	// Use registered service to do everything related
	// to account update.
	ves, err := ep.Service.UpdateAccount(s, &account)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.Log.Debug("Validation errors", "dump", fmt.Sprintf("%+v", ves.FieldErrors))
		ep.rerenderAccountForm(w, r, account.ToForm(), ves, kbs.NewTmpl, accountCreateAction())
		return
	}

	// Non validation errors
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), UpdateAccountErrMsg, err)
		return
	}

	m := ep.Localize(r, AccountUpdatedInfoMsg)
	ep.RedirectWithFlash(w, r, AccountPath(), m, kbs.InfoMT)
}

// InitDeleteAccount web endpoint.
func (ep *Endpoint) InitDeleteAccount(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the account from repo.
	account, err := ep.Service.GetAccount(s)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), GetAccountsErrMsg, err)
		return
	}

	// Wrap response
	accountForm := account.ToForm()
	wr := ep.WrapRes(w, r, &accountForm, nil)
	wr.SetAction(accountDeleteAction(&accountForm))

	// Template
	ts, err := ep.TemplateFor(accountRes, kbs.InitDelTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), CannotProcErrMsg, err)
		return
	}
}

// DeleteAccount web endpoint.
func (ep *Endpoint) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), DeleteAccountErrMsg, err)
		return
	}

	// Service
	err = ep.Service.DeleteAccount(s)
	if err != nil {
		ep.ErrorRedirect(w, r, AccountPath(), DeleteAccountErrMsg, err)
		return
	}

	m := ep.Localize(r, AccountDeletedInfoMsg)
	ep.RedirectWithFlash(w, r, AccountPath(), m, kbs.InfoMT)
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
