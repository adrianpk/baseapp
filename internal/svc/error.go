package svc

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
)

type (
	Err struct {
		kbs.Err
	}

	noRepoError struct {
		kbs.Err
	}
)

var (
	NoRepoErr           = NewErr("no_repo_err", nil)
	AlreadyConfirmedErr = NewErr("already_confirmed_err_msg", nil)
	CredentialsErr      = NewErr("creadentials_err_msg", nil)
)

func NewErr(msgID string, err error) Err {
	return Err{
		kbs.NewErr(msgID, err),
	}
}
