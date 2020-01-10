package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	AuthRepo interface {
		CreateAccountRole(u *model.AccountRole, tx ...*sqlx.Tx) error
		GetAllAccountRoles() (auth []model.AccountRole, err error)
		GetAccountRole(id uuid.UUID) (auth model.AccountRole, err error)
		GetAccountRoleBySlug(slug string) (auth model.AccountRole, err error)
		GetAccountRoleByAccountID(uuid.UUID) ([]model.AccountRole, error)
		GetAccountRoleByRoleID(uuid.UUID) ([]model.AccountRole, error)
		UpdateAccountRole(auth *model.AccountRole, tx ...*sqlx.Tx) error
		DeleteAccountRole(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteAccountRoleBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
