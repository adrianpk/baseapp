package svc

import (
	"errors"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	PermissionValidator struct {
		Model *model.Permission
		kbs.Validator
	}
)

func NewPermissionValidator(u *model.Permission) PermissionValidator {
	return PermissionValidator{
		Model:     u,
		Validator: kbs.NewValidator(),
	}
}

func (uv PermissionValidator) ValidateForCreate() error {
	// Name
	ok0 := uv.ValidateRequiredName()

	if ok0 {
		return nil
	}

	return errors.New("permission has errors")
}

// NOTE: Update validations shoud be different
// than the ones executed on creation.
func (uv PermissionValidator) ValidateForUpdate() error {
	// Name
	ok0 := uv.ValidateRequiredName()

	if ok0 {
		return nil
	}

	return errors.New("permission has errors")
}

func (uv PermissionValidator) ValidateRequiredName(errMsg ...string) (ok bool) {
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
