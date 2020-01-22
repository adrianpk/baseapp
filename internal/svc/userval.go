package svc

import (
	"errors"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	UserValidator struct {
		Model *model.User
		kbs.Validator
	}
)

func NewUserValidator(u *model.User) UserValidator {
	return UserValidator{
		Model:     u,
		Validator: kbs.NewValidator(),
	}
}

func (uv UserValidator) ValidateForCreate() error {
	// Username
	ok0 := uv.ValidateRequiredUsername()
	ok1 := uv.ValidateMinLengthUsername(4)
	ok2 := uv.ValidateMaxLengthUsername(16)
	// Email
	ok3 := uv.ValidateEmailEmail()
	ok4 := uv.ValidateEmailConfirmation()
	// Password
	ok5 := uv.ValidateRequiredPassword()
	ok6 := uv.ValidateMinLengthPassword(8)
	ok7 := uv.ValidateMaxLengthPassword(32)

	if ok0 && ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7 {
		return nil
	}

	return errors.New("user has errors")
}

// NOTE: Update validations shoud be different
// than the ones executed on creation.
func (uv UserValidator) ValidateForUpdate() error {
	// Username
	ok0 := uv.ValidateRequiredUsername()
	ok1 := uv.ValidateMinLengthUsername(4)
	ok2 := uv.ValidateMaxLengthUsername(16)
	// Email
	ok3 := uv.ValidateEmailEmail()
	ok4 := uv.ValidateEmailConfirmation()
	// Password
	//ok5 := uv.ValidateRequiredPassword()
	//ok6 := uv.ValidateMinLengthPassword(8)
	//ok7 := uv.ValidateMaxLengthPassword(32)

	if ok0 && ok1 && ok2 && ok3 && ok4 {
		return nil
	}

	return errors.New("user has errors")
}

func (uv UserValidator) ValidateForSignUp() error {
	// Username
	ok0 := uv.ValidateRequiredUsername()
	ok1 := uv.ValidateMinLengthUsername(4)
	ok2 := uv.ValidateMaxLengthUsername(16)
	// Email
	ok3 := uv.ValidateEmailEmail()
	ok4 := uv.ValidateEmailConfirmation()
	// Password
	ok5 := uv.ValidateRequiredPassword()
	ok6 := uv.ValidateMinLengthPassword(8)
	ok7 := uv.ValidateMaxLengthPassword(32)

	if ok0 && ok1 && ok2 && ok3 && ok4 && ok5 && ok6 && ok7 {
		return nil
	}

	return errors.New("user has errors")
}

func (uv UserValidator) ValidateRequiredUsername(errMsg ...string) (ok bool) {
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

func (uv UserValidator) ValidateMinLengthUsername(min int, errMsg ...string) (ok bool) {
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

func (uv UserValidator) ValidateMaxLengthUsername(max int, errMsg ...string) (ok bool) {
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

func (uv UserValidator) ValidateEmailEmail(errMsg ...string) (ok bool) {
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

func (uv UserValidator) ValidateEmailConfirmation(errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateConfirmation(u.Email.String, u.EmailConfirmation.String)
	if ok {
		return true
	}

	msg := kbs.ValMsg.NoMatchErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Email"] = append(uv.Errors["Email"], msg)
	uv.Errors["EmailConfirmation"] = append(uv.Errors["EmailConfirmation"], msg)
	return false
}

func (uv UserValidator) ValidateRequiredPassword(errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateRequired(u.Password)
	if ok {
		return true
	}

	msg := kbs.ValMsg.RequiredErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Password"] = append(uv.Errors["Password"], msg)
	return false
}

func (uv UserValidator) ValidateMinLengthPassword(min int, errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateMinLength(u.Password, min)
	if ok {
		return true
	}

	msg := kbs.ValMsg.MinLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Password"] = append(uv.Errors["Password"], msg)
	return false
}

func (uv UserValidator) ValidateMaxLengthPassword(max int, errMsg ...string) (ok bool) {
	u := uv.Model

	ok = uv.ValidateMaxLength(u.Password, max)
	if ok {
		return true
	}

	msg := kbs.ValMsg.MaxLengthErrMsg
	if len(errMsg) > 0 {
		msg = errMsg[0]
	}

	uv.Errors["Password"] = append(uv.Errors["Password"], msg)
	return false
}
