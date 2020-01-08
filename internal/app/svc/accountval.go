package svc

import (
	"errors"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	AccountValidator struct {
		Model *model.Account
		kbs.Validator
	}
)

func NewAccountValidator(u *model.Account) AccountValidator {
	return AccountValidator{
		Model:     u,
		Validator: kbs.NewValidator(),
	}
}

func (uv AccountValidator) ValidateForCreate() error {
	// Accountname
	ok0 := uv.ValidateRequiredName()
	ok1 := uv.ValidateMinLengthName(4)
	ok2 := uv.ValidateMaxLengthName(16)
	// Email
	ok3 := uv.ValidateEmailEmail()

	if ok0 && ok1 && ok2 && ok3 {
		return nil
	}

	return errors.New("account has errors")
}

// NOTE: Update validations shoud be different
// than the ones executed on creation.
func (uv AccountValidator) ValidateForUpdate() error {
	// Name
	ok0 := uv.ValidateRequiredName()
	ok1 := uv.ValidateMinLengthName(4)
	ok2 := uv.ValidateMaxLengthName(16)
	// Email
	ok3 := uv.ValidateEmailEmail()

	if ok0 && ok1 && ok2 && ok3 {
		return nil
	}

	return errors.New("account has errors")
}

func (uv AccountValidator) ValidateRequiredName(errMsg ...string) (ok bool) {
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

func (uv AccountValidator) ValidateMinLengthName(min int, errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateMinLength(u.Name.String, min)
	if ok {
		return true
	}

	msg := kbs.ValMsg.MinLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Name"] = append(uv.Errors["Name"], msg)
	return false
}

func (uv AccountValidator) ValidateMaxLengthName(max int, errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateMaxLength(u.Name.String, max)
	if ok {
		return true
	}

	msg := kbs.ValMsg.MaxLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Name"] = append(uv.Errors["Name"], msg)
	return false
}

func (uv AccountValidator) ValidateEmailEmail(errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateEmail(u.Email.String)
	if ok {
		return true
	}

	msg := kbs.ValMsg.NotEmailErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Email"] = append(uv.Errors["Email"], msg)
	return false
}
