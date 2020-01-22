package svc

import (
	"errors"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	ResourceValidator struct {
		Model *model.Resource
		kbs.Validator
	}
)

func NewResourceValidator(u *model.Resource) ResourceValidator {
	return ResourceValidator{
		Model:     u,
		Validator: kbs.NewValidator(),
	}
}

func (uv ResourceValidator) ValidateForCreate() error {
	// Name
	ok0 := uv.ValidateRequiredName()
	ok1 := uv.ValidateRequiredPath()

	if ok0 && ok1 {
		return nil
	}

	return errors.New("resource has errors")
}

// NOTE: Update validations shoud be different
// than the ones executed on creation.
func (uv ResourceValidator) ValidateForUpdate() error {
	// Name
	ok0 := uv.ValidateRequiredName()
	ok1 := uv.ValidateRequiredPath()

	if ok0 && ok1 {
		return nil
	}

	return errors.New("resource has errors")
}

func (uv ResourceValidator) ValidateRequiredName(errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateRequired(u.Name.String)
	if ok {
		return true
	}

	msg := kbs.ValMsg.RequiredErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Name"] = append(uv.Errors["Name"], msg)
	return false
}

func (uv ResourceValidator) ValidateRequiredPath(errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateRequired(u.Path.String)
	if ok {
		return true
	}

	msg := kbs.ValMsg.RequiredErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Path"] = append(uv.Errors["Path"], msg)
	return false
}
