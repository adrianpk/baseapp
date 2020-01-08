package svc

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/repo"
	//repo "gitlab.com/kabestan/repo/baseapp/internal/repo/pg"
)

type (
	Service struct {
		*kbs.Service
		UserRepo    repo.UserRepo
		AccountRepo repo.AccountRepo
	}
)

func NewService(cfg *kbs.Config, log kbs.Logger, name string) *Service {
	return &Service{
		Service: kbs.NewService(cfg, log, name),
	}
}
