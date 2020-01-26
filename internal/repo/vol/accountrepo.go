package pg

import (
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

// Custom

func (ar *AuthRepo) GetAccountRoles(slug string) (roles []model.Role, err error) {
	panic("not implemented")
}

func (ar *AuthRepo) GetNotAccountRoles(slug string) (roles []model.Role, err error) {
	panic("not implemented")
}

func (ar *AuthRepo) AppendAccountRole(accountSlug, roleSlug string) (err error) {
	panic("not implemented")
}

func (ar *AuthRepo) RemoveAccountRole(accountSlug, roleSlug string) (err error) {
	panic("not implemented")
}
