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
	ok0 := uv.ValidateRequiredUsername()
	ok1 := uv.ValidateMinLengthUsername(4)
	ok2 := uv.ValidateMaxLengthUsername(16)
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
	// Username
	ok0 := uv.ValidateRequiredUsername()
	ok1 := uv.ValidateMinLengthUsername(4)
	ok2 := uv.ValidateMaxLengthUsername(16)
	// Email
	ok3 := uv.ValidateEmailEmail()

	if ok0 && ok1 && ok2 && ok3 {
		return nil
	}

	return errors.New("account has errors")
}

func (uv AccountValidator) ValidateRequiredUsername(errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateRequired(u.Username.String)
	if ok {
		return true
	}

	msg := kbs.ValMsg.RequiredErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Username"] = append(uv.Errors["Username"], msg)
	return false
}

func (uv AccountValidator) ValidateMinLengthUsername(min int, errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateMinLength(u.Username.String, min)
	if ok {
		return true
	}

	msg := kbs.ValMsg.MinLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Username"] = append(uv.Errors["Username"], msg)
	return false
}

func (uv AccountValidator) ValidateMaxLengthUsername(max int, errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateMaxLength(u.Username.String, max)
	if ok {
		return true
	}

	msg := kbs.ValMsg.MaxLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Username"] = append(uv.Errors["Username"], msg)
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
