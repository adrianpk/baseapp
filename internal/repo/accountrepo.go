package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	AccountRepo interface {
		Create(u *model.Account, tx ...*sqlx.Tx) error
		GetAll() (accounts []model.Account, err error)
		Get(id uuid.UUID) (account model.Account, err error)
		GetBySlug(slug string) (account model.Account, err error)
		GetByOwnerID(id uuid.UUID) (account model.Account, err error)
		GetByName(name string) (model.Account, error)
		Update(account *model.Account, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
