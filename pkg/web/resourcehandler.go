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
	ResourceCreatedInfoMsg = "resource_created_info_msg"
	ResourceUpdatedInfoMsg = "resource_updated_info_msg"
	ResourceDeletedInfoMsg = "resource_deleted_info_msg"
	// Error
	CreateResourceErrMsg = "create_resource_err_msg"
	IndexResourcesErrMsg = "get_all_resources_err_msg"
	GetResourceErrMsg    = "get_resource_err_msg"
	GetResourcesErrMsg   = "get_resources_err_msg"
	UpdateResourceErrMsg = "update_resource_err_msg"
	DeleteResourceErrMsg = "delete_resource_err_msg"
)

// IndexResources web endpoint.
func (ep *Endpoint) IndexResources(w http.ResponseWriter, r *http.Request) {
	// Get resources list from registered service
	resources, err := ep.Service.IndexResources()
	if err != nil {
		ep.ErrorRedirect(w, r, "/", IndexResourcesErrMsg, err)
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
	m := ep.Localize(r, ResourceCreatedInfoMsg)
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
		ep.ErrorRedirect(w, r, ResourcePath(), GetResourceErrMsg, err)
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
		ep.ErrorRedirect(w, r, ResourcePath(), GetResourceErrMsg, err)
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
		ep.ErrorRedirect(w, r, ResourcePath(), GetResourceErrMsg, err)
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
		ep.ErrorRedirect(w, r, ResourcePath(), UpdateResourceErrMsg, err)
		return
	}

	m := ep.Localize(r, ResourceUpdatedInfoMsg)
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
		ep.ErrorRedirect(w, r, ResourcePath(), GetResourcesErrMsg, err)
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
		ep.ErrorRedirect(w, r, ResourcePath(), DeleteResourceErrMsg, err)
		return
	}

	// Service
	err = ep.Service.DeleteResource(s)
	if err != nil {
		ep.ErrorRedirect(w, r, ResourcePath(), DeleteResourceErrMsg, err)
		return
	}

	m := ep.Localize(r, ResourceDeletedInfoMsg)
	ep.RedirectWithFlash(w, r, ResourcePath(), m, kbs.InfoMT)
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

func (ep *Endpoint) ErrorRedirect(w http.ResponseWriter, r *http.Request, redirPath, msgID string, err error) {
	m := ep.Localize(r, msgID)
	ep.RedirectWithFlash(w, r, redirPath, m, kbs.ErrorMT)
	ep.Log.Error(err)
}