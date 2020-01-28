package pg

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
)

type (
	AuthRepo struct {
		Cfg  *kbs.Config
		Log  kbs.Logger
		Name string
	}
)

func NewAuthRepo(cfg *kbs.Config, log kbs.Logger, name string) *AuthRepo {
	return &AuthRepo{
		Cfg:  cfg,
		Log:  log,
		Name: name,
	}
}

func (ar *AuthRepo) GetResourcePermissionTagsByPath(path string) (tags []string, err error) {
	panic("not implemented")
}
