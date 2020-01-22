package svc

import (
	"errors"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	RoleValidator struct {
		Model *model.Role
		kbs.Validator
	}
)

func NewRoleValidator(u *model.Role) RoleValidator {
	return RoleValidator{
		Model:     u,
		Validator: kbs.NewValidator(),
	}
}

func (uv RoleValidator) ValidateForCreate() error {
	// Name
	ok0 := uv.ValidateRequiredName()

	if ok0 {
		return nil
	}

	return errors.New("role has errors")
}

// NOTE: Update validations shoud be different
// than the ones executed on creation.
func (uv RoleValidator) ValidateForUpdate() error {
	// Name
	ok0 := uv.ValidateRequiredName()

	if ok0 {
		return nil
	}

	return errors.New("role has errors")
}

func (uv RoleValidator) ValidateRequiredName(errMsg ...string) (ok bool) {
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
