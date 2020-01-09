package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	RoleRepo interface {
		Create(u *model.Role, tx ...*sqlx.Tx) error
		GetAll() (roles []model.Role, err error)
		Get(id uuid.UUID) (role model.Role, err error)
		GetBySlug(slug string) (role model.Role, err error)
		GetByName(name string) (model.Role, error)
		Update(role *model.Role, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
