package web

import (
	"fmt"
	"net/http"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

const (
	permissionRes = "permission"
)

const (
	// Defined in 'assets/web/embed/i18n/xx.json'
	permissionResPl = "permissions"
)

// IndexPermissions web endpoint.
func (ep *Endpoint) IndexPermissions(w http.ResponseWriter, r *http.Request) {
	// Get permissions list from registered service
	permissions, err := ep.Service.IndexPermissions()
	if err != nil {
		ep.ErrorRedirect(w, r, "/", IndexErrMsg, err)
		return
	}

	// Convert result list into a form list
	// Models use sql null types but templates looks
	// clearer if we use plain Go type.
	// i.e.: $permission.Permissionname instead of $permission.Permissionname.String
	l := model.ToPermissionFormList(permissions)
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(permissionRes, kbs.IndexTmpl)
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

func (ep *Endpoint) NewPermission(w http.ResponseWriter, r *http.Request) {
	permissionForm := model.PermissionForm{}

	// Wrap response
	wr := ep.WrapRes(w, r, &permissionForm, nil)
	wr.SetAction(permissionCreateAction())

	// Get template to render from cache.
	ts, err := ep.TemplateFor(permissionRes, kbs.NewTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}
}

func (ep *Endpoint) CreatePermission(w http.ResponseWriter, r *http.Request) {
	// Decode request data into a form.
	permissionForm := model.PermissionForm{}
	err := ep.FormToModel(r, &permissionForm)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	permission := permissionForm.ToModel()

	// Use registered service to do everything related
	// to permission creation.
	ves, err := ep.Service.CreatePermission(&permission)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.rerenderPermissionForm(w, r, permission.ToForm(), ves, kbs.NewTmpl, permissionCreateAction())
		return
	}

	// Then take care of other kind of possible errors
	// that service can generate.
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, CreatedInfoMsg)
	ep.RedirectWithFlash(w, r, PermissionPath(), m, kbs.InfoMT)
}

// ShowPermission web endpoint.
func (ep *Endpoint) ShowPermission(w http.ResponseWriter, r *http.Request) {
	// Get slug from request context.
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Use registered service to do everything related
	// to permission creation.
	permission, err := ep.Service.GetPermission(s)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	// Wrap response
	wr := ep.WrapRes(w, r, permission.ToForm(), nil)

	// Template
	ts, err := ep.TemplateFor(permissionRes, kbs.ShowTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}
}

// EditPermission web endpoint.
func (ep *Endpoint) EditPermission(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the permission from repo.
	permission, err := ep.Service.GetPermission(s)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	// Wrap response
	permissionForm := permission.ToForm()
	wr := ep.WrapRes(w, r, &permissionForm, nil)
	wr.SetAction(permissionUpdateAction(&permissionForm))

	// Template
	ts, err := ep.TemplateFor(permissionRes, kbs.EditTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}
}

// UpdatePermission web endpoint.
func (ep *Endpoint) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	// Decode request data into a form.
	permissionForm := model.PermissionForm{}
	err = ep.FormToModel(r, &permissionForm)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	permission := permissionForm.ToModel()

	// Use registered service to do everything related
	// to permission update.
	ves, err := ep.Service.UpdatePermission(s, &permission)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.Log.Debug("Validation errors", "dump", fmt.Sprintf("%+v", ves.FieldErrors))
		ep.rerenderPermissionForm(w, r, permission.ToForm(), ves, kbs.NewTmpl, permissionCreateAction())
		return
	}

	// Non validation errors
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	m := ep.Localize(r, UpdatedInfoMsg)
	ep.RedirectWithFlash(w, r, PermissionPath(), m, kbs.InfoMT)
}

// InitDeletePermission web endpoint.
func (ep *Endpoint) InitDeletePermission(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the permission from repo.
	permission, err := ep.Service.GetPermission(s)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	// Wrap response
	permissionForm := permission.ToForm()
	wr := ep.WrapRes(w, r, &permissionForm, nil)
	wr.SetAction(permissionDeleteAction(&permissionForm))

	// Template
	ts, err := ep.TemplateFor(permissionRes, kbs.InitDelTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}
}

// DeletePermission web endpoint.
func (ep *Endpoint) DeletePermission(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	// Service
	err = ep.Service.DeletePermission(s)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), GetErrMsg, err)
		return
	}

	m := ep.Localize(r, DeletedInfoMsg)
	ep.RedirectWithFlash(w, r, PermissionPath(), m, kbs.InfoMT)
}

func (ep *Endpoint) rerenderPermissionForm(w http.ResponseWriter, r *http.Request, data interface{}, valErrors kbs.ValErrorSet, template string, action kbs.FormAction) {
	wr := ep.WrapRes(w, r, data, valErrors)
	wr.AddErrorFlash(InputValuesErrMsg)
	wr.SetAction(action)

	ts, err := ep.TemplateFor(permissionRes, template)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), InputValuesErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, PermissionPath(), CannotProcErrMsg, err)
		return
	}

	return
}

// permissionCreateAction
func permissionCreateAction() kbs.FormAction {
	return kbs.FormAction{Target: fmt.Sprintf("%s", PermissionPath()), Method: "POST"}
}

// permissionUpdateAction
func permissionUpdateAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: PermissionPathSlug(model), Method: "PUT"}
}

// permissionDeleteAction
func permissionDeleteAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: PermissionPathSlug(model), Method: "DELETE"}
}
