package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	AccountRoleRepo interface {
		Create(u *model.AccountRole, tx ...*sqlx.Tx) error
		GetAll() (accountRoles []model.AccountRole, err error)
		Get(id uuid.UUID) (accountRole model.AccountRole, err error)
		GetBySlug(slug string) (accountRole model.AccountRole, err error)
		GetByAccountID(uuid.UUID) ([]model.AccountRole, error)
		GetByRoleID(uuid.UUID) ([]model.AccountRole, error)
		Update(accountRole *model.AccountRole, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
