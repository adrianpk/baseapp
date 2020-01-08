package svc

import (
	"github.com/jmoiron/sqlx"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/repo"
	//repo "gitlab.com/kabestan/repo/baseapp/internal/repo/pg"
)

type (
	Service struct {
		*kbs.Service
		DB          *sqlx.DB
		UserRepo    repo.UserRepo
		AccountRepo repo.AccountRepo
	}
)

func NewService(cfg *kbs.Config, log kbs.Logger, name string, db *sqlx.DB) *Service {
	return &Service{
		Service: kbs.NewService(cfg, log, name),
		DB:      db,
	}
}

func (s *Service) getTx() (tx *sqlx.Tx, err error) {
	tx, err = s.DB.Beginx()
	if err != nil {
		return tx, err
	}

	return tx, nil
}
