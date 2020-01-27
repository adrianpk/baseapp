package web

import (
	"fmt"
	"net/http"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

const (
	resourceRes = "resource"
)

const (
	// Defined in 'assets/web/embed/i18n/xx.json'
	resourceResPl = "resources"
)

// IndexResources web endpoint.
func (ep *Endpoint) IndexResources(w http.ResponseWriter, r *http.Request) {
	// Get resources list from registered service
	resources, err := ep.Service.IndexResources()
	if err != nil {
		ep.ErrorRedirect(w, r, "/", IndexErrMsg, err)
		return
	}

	// Convert result list into a form list
	// Models use sql null types but templates looks
	// clearer if we use plain Go type.
	// i.e.: $resource.Resourcename instead of $resource.Resourcename.String
	l := model.ToResourceFormList(resources)
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(resourceRes, kbs.IndexTmpl)
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

func (ep *Endpoint) NewResource(w http.ResponseWriter, r *http.Request) {
	resourceForm := model.ResourceForm{}

	// Wrap response
	wr := ep.WrapRes(w, r, &resourceForm, nil)
	wr.SetAction(resourceCreateAction())

	// Get template to render from cache.
	ts, err := ep.TemplateFor(resourceRes, kbs.NewTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	// Execute it and redirect if error.
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}
}

func (ep *Endpoint) CreateResource(w http.ResponseWriter, r *http.Request) {
	// Decode request data into a form.
	resourceForm := model.ResourceForm{}
	err := ep.FormToModel(r, &resourceForm)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	resource := resourceForm.ToModel()

	// Use registered service to do everything related
	// to resource creation.
	ves, err := ep.Service.CreateResource(&resource)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.rerenderResourceForm(w, r, resource.ToForm(), ves, kbs.NewTmpl, resourceCreateAction())
		return
	}

	// Then take care of other kind of possible errors
	// that service can generate.
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, CreatedInfoMsg)
	ep.RedirectWithFlash(w, r, ResourcePath(), m, kbs.InfoMT)
}

// ShowResource web endpoint.
func (ep *Endpoint) ShowResource(w http.ResponseWriter, r *http.Request) {
	// Get slug from request context.
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Use registered service to do everything related
	// to resource creation.
	resource, err := ep.Service.GetResource(s)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	// Wrap response
	wr := ep.WrapRes(w, r, resource.ToForm(), nil)

	// Template
	ts, err := ep.TemplateFor(resourceRes, kbs.ShowTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}
}

// EditResource web endpoint.
func (ep *Endpoint) EditResource(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the resource from repo.
	resource, err := ep.Service.GetResource(s)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	// Wrap response
	resourceForm := resource.ToForm()
	wr := ep.WrapRes(w, r, &resourceForm, nil)
	wr.SetAction(resourceUpdateAction(&resourceForm))

	// Template
	ts, err := ep.TemplateFor(resourceRes, kbs.EditTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}
}

// UpdateResource web endpoint.
func (ep *Endpoint) UpdateResource(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	// Decode request data into a form.
	resourceForm := model.ResourceForm{}
	err = ep.FormToModel(r, &resourceForm)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Create a model using form values.
	resource := resourceForm.ToModel()

	// Use registered service to do everything related
	// to resource update.
	ves, err := ep.Service.UpdateResource(s, &resource)

	// First take care of service validation errors.
	if !ves.IsEmpty() {
		ep.Log.Debug("Validation errors", "dump", fmt.Sprintf("%+v", ves.FieldErrors))
		ep.rerenderResourceForm(w, r, resource.ToForm(), ves, kbs.NewTmpl, resourceCreateAction())
		return
	}

	// Non validation errors
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	m := ep.Localize(r, UpdatedInfoMsg)
	ep.RedirectWithFlash(w, r, ResourcePath(), m, kbs.InfoMT)
}

// InitDeleteResource web endpoint.
func (ep *Endpoint) InitDeleteResource(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get the resource from repo.
	resource, err := ep.Service.GetResource(s)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	// Wrap response
	resourceForm := resource.ToForm()
	wr := ep.WrapRes(w, r, &resourceForm, nil)
	wr.SetAction(resourceDeleteAction(&resourceForm))

	// Template
	ts, err := ep.TemplateFor(resourceRes, kbs.InitDelTmpl)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}
}

// DeleteResource web endpoint.
func (ep *Endpoint) DeleteResource(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	// Service
	err = ep.Service.DeleteResource(s)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), GetErrMsg, err)
		return
	}

	m := ep.Localize(r, DeletedInfoMsg)
	ep.RedirectWithFlash(w, r, ResourcePath(), m, kbs.InfoMT)
}

// IndexResourcePermissions web endpoint.
func (ep *Endpoint) IndexResourcePermissions(w http.ResponseWriter, r *http.Request) {
	s, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to get resource.
	resource, err := ep.Service.GetResource(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all available permissions from repo.
	notApplied, err := ep.Service.GetNotResourcePermissions(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Use registerd service to get all resource associated permissions from repo.
	applied, err := ep.Service.GetResourcePermissions(s)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}

	// Group both lists into a map
	l := map[string]interface{}{
		"resource":    resource.ToForm(),
		"not-applied": model.ToPermissionFormList(notApplied),
		"applied":     model.ToPermissionFormList(applied),
	}

	// Wrap response
	wr := ep.WrapRes(w, r, l, nil)

	// Get template to render from cache.
	ts, err := ep.TemplateFor(resourceRes, PermissionsTmpl)
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

// AppendResourcePermission web endpoint.
func (ep *Endpoint) AppendResourcePermission(w http.ResponseWriter, r *http.Request) {
	resourceSlug, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	resourceForm := model.ResourceForm{Slug: resourceSlug}

	// Decode request data into a form.
	permissionForm := model.PermissionForm{}
	err = ep.FormToModel(r, &permissionForm)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePathPermissions(resourceForm), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to append permission.
	err = ep.Service.AppendResourcePermission(resourceSlug, permissionForm.Slug)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}
	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, PermissionAppendedInfoMsg)
	ep.RedirectWithFlash(w, r, ResourcePathPermissions(resourceForm), m, kbs.InfoMT)
}

// RemoveResourcePermission web endpoint.
func (ep *Endpoint) RemoveResourcePermission(w http.ResponseWriter, r *http.Request) {
	resourceSlug, err := ep.getSlug(r)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), CannotProcErrMsg, err)
		return
	}

	resourceForm := model.ResourceForm{Slug: resourceSlug}

	// Decode request data into a form.
	permissionForm := model.PermissionForm{}
	err = ep.FormToModel(r, &permissionForm)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePathPermissions(resourceForm), CannotProcErrMsg, err)
		return
	}

	// Use registerd service to append permission.
	err = ep.Service.RemoveResourcePermission(resourceSlug, permissionForm.Slug)
	if err != nil {
		ep.ErrorRedirect(w, r, UserPath(), GetErrMsg, err)
		return
	}
	// Localize Ok info message, put it into a flash message
	// and redirect to index.
	m := ep.Localize(r, PermissionRemovedInfoMsg)

	ep.RedirectWithFlash(w, r, ResourcePathPermissions(resourceForm), m, kbs.InfoMT)
}

func (ep *Endpoint) rerenderResourceForm(w http.ResponseWriter, r *http.Request, data interface{}, valErrors kbs.ValErrorSet, template string, action kbs.FormAction) {
	wr := ep.WrapRes(w, r, data, valErrors)
	wr.AddErrorFlash(InputValuesErrMsg)
	wr.SetAction(action)

	ts, err := ep.TemplateFor(resourceRes, template)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), InputValuesErrMsg, err)
		return
	}

	// Write response
	err = ts.Execute(w, wr)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), CannotProcErrMsg, err)
		return
	}

	return
}

// resourceCreateAction
func resourceCreateAction() kbs.FormAction {
	return kbs.FormAction{Target: fmt.Sprintf("%s", ResourcePath()), Method: "POST"}
}

// resourceUpdateAction
func resourceUpdateAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: ResourcePathSlug(model), Method: "PUT"}
}

// resourceDeleteAction
func resourceDeleteAction(model kbs.Identifiable) kbs.FormAction {
	return kbs.FormAction{Target: ResourcePathSlug(model), Method: "DELETE"}
}
