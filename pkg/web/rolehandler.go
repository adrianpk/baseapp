package web

import (
	"fmt"
	"net/http"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

const (
	roleRes = "role"
)

const (
	// Defined in 'assets/web/embed/i18n/xx.json'
	roleResPl = "roles"
)

const (
	PermissionsTmpl = "permissions.tmpl"
	// Info
	PermissionAppendedInfoMsg = "permission_appended_info_msg"
	PermissionRemovedInfoMsg  = "permission_removed_info_msg"
)

// IndexRoles web endpoint.
func (ep *Endpoint) IndexRoles(w http.ResponseWriter, r *http.Request) {
	// Get roles list from registered service
	roles, err := ep.Service().IndexRoles()
	if err != nil {
		ep.ErrorRedirect(w, r, "/", IndexErrMsg, err)
		return
	}

	// Convert result list into a form list
	// Models use sql null types but templates looks
	// clearer if we use plain Go type.
	// i.e.: $role.Rolename instead of $role.Rolename.String
	l := model.ToRoleFormList(roles)
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(roleRes, kbs.IndexTmpl)
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

func (ep *Endpoint) NewRole(w http.ResponseWriter, r *http.Request) {
	roleForm := model.RoleForm{}

	// Wrap response
	wr := ep.WrapRes(w, r, &roleForm, nil)
	wr.SetAction(roleCreateAction())

	// Get template to render from cache.
	ts, err := ep.TemplateFor(roleRes, kbs.NewTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}
}

func (ep *Endpoint) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Decode request data into a form.
	roleForm := model.RoleForm{}
	err := ep.FormToModel(r, &roleForm)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	role := roleForm.ToModel()

	// Use registered service to do everything related
	// to role creation.
	ves, err := ep.Service().CreateRole(&role)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.rerenderRoleForm(w, r, role.ToForm(), ves, kbs.NewTmpl, roleCreateAction())
		return
	}

	// Then take care of other kind of possible errors
	// that service can generate.
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, CreatedInfoMsg)
	ep.RedirectWithFlash(w, r, RolePath(), m, kbs.InfoMT)
}

// ShowRole web endpoint.
func (ep *Endpoint) ShowRole(w http.ResponseWriter, r *http.Request) {
	// Get slug from request context.
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Use registered service to do everything related
	// to role creation.
	role, err := ep.Service().GetRole(s)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	// Wrap response
	wr := ep.WrapRes(w, r, role.ToForm(), nil)

	// Template
	ts, err := ep.TemplateFor(roleRes, kbs.ShowTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}
}

// EditRole web endpoint.
func (ep *Endpoint) EditRole(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the role from repo.
	role, err := ep.Service().GetRole(s)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	// Wrap response
	roleForm := role.ToForm()
	wr := ep.WrapRes(w, r, &roleForm, nil)
	wr.SetAction(roleUpdateAction(&roleForm))

	// Template
	ts, err := ep.TemplateFor(roleRes, kbs.EditTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}
}

// UpdateRole web endpoint.
func (ep *Endpoint) UpdateRole(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	// Decode request data into a form.
	roleForm := model.RoleForm{}
	err = ep.FormToModel(r, &roleForm)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	role := roleForm.ToModel()

	// Use registered service to do everything related
	// to role update.
	ves, err := ep.Service().UpdateRole(s, &role)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.Log.Debug("Validation errors", "dump", fmt.Sprintf("%+v", ves.FieldErrors))
		ep.rerenderRoleForm(w, r, role.ToForm(), ves, kbs.NewTmpl, roleCreateAction())
		return
	}

	// Non validation errors
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	m := ep.Localize(r, UpdatedInfoMsg)
	ep.RedirectWithFlash(w, r, RolePath(), m, kbs.InfoMT)
}

// InitDeleteRole web endpoint.
func (ep *Endpoint) InitDeleteRole(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the role from repo.
	role, err := ep.Service().GetRole(s)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	// Wrap response
	roleForm := role.ToForm()
	wr := ep.WrapRes(w, r, &roleForm, nil)
	wr.SetAction(roleDeleteAction(&roleForm))

	// Template
	ts, err := ep.TemplateFor(roleRes, kbs.InitDelTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}
}

// DeleteRole web endpoint.
func (ep *Endpoint) DeleteRole(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	// Service
	err = ep.Service().DeleteRole(s)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), GetErrMsg, err)
		return
	}

	m := ep.Localize(r, DeletedInfoMsg)
	ep.RedirectWithFlash(w, r, RolePath(), m, kbs.InfoMT)
}

// IndexRolePermissions web endpoint.
func (ep *Endpoint) IndexRolePermissions(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get role.
	role, err := ep.Service().GetRole(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all available permissions from repo.
	notApplied, err := ep.Service().GetNotRolePermissions(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all role associated permissions from repo.
	applied, err := ep.Service().GetRolePermissions(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Group both lists into a map
	l := map[string]interface{}{
		"role":        role.ToForm(),
		"not-applied": model.ToPermissionFormList(notApplied),
		"applied":     model.ToPermissionFormList(applied),
	}

	// Wrap response
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(roleRes, PermissionsTmpl)
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

// AppendRolePermission web endpoint.
func (ep *Endpoint) AppendRolePermission(w http.ResponseWriter, r *http.Request) {
	roleSlug, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	roleForm := model.RoleForm{Slug: roleSlug}

	// Decode request data into a form.
	permissionForm := model.PermissionForm{}
	err = ep.FormToModel(r, &permissionForm)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePathPermissions(roleForm), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to append permission.
	err = ep.Service().AppendRolePermission(roleSlug, permissionForm.Slug)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}
	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, PermissionAppendedInfoMsg)
	ep.RedirectWithFlash(w, r, RolePathPermissions(roleForm), m, kbs.InfoMT)
}

// RemoveRolePermission web endpoint.
func (ep *Endpoint) RemoveRolePermission(w http.ResponseWriter, r *http.Request) {
	roleSlug, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	roleForm := model.RoleForm{Slug: roleSlug}

	// Decode request data into a form.
	permissionForm := model.PermissionForm{}
	err = ep.FormToModel(r, &permissionForm)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePathPermissions(roleForm), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to append permission.
	err = ep.Service().RemoveRolePermission(roleSlug, permissionForm.Slug)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}
	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, PermissionRemovedInfoMsg)

	ep.RedirectWithFlash(w, r, RolePathPermissions(roleForm), m, kbs.InfoMT)
}

func (ep *Endpoint) rerenderRoleForm(w http.ResponseWriter, r *http.Request, data interface{}, valErrors kbs.ValErrorSet, template string, action kbs.FormAction) {
	wr := ep.WrapRes(w, r, data, valErrors)
	wr.AddErrorFlash(InputValuesErrMsg)
	wr.SetAction(action)

	ts, err := ep.TemplateFor(roleRes, template)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), InputValuesErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, RolePath(), CannotProcErrMsg, err)
		return
	}

	return
}

// roleCreateAction
func roleCreateAction() kbs.FormAction {
	return kbs.FormAction{Target: fmt.Sprintf("%s", RolePath()), Method: "POST"}
}

// roleUpdateAction
func roleUpdateAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: RolePathSlug(model), Method: "PUT"}
}

// roleDeleteAction
func roleDeleteAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: RolePathSlug(model), Method: "DELETE"}
}
