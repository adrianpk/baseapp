package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	PermissionRepo interface {
		Create(u *model.Permission, tx ...*sqlx.Tx) error
		GetAll() (permissions []model.Permission, err error)
		Get(id uuid.UUID) (permission model.Permission, err error)
		GetBySlug(slug string) (permission model.Permission, err error)
		GetByName(name string) (model.Permission, error)
		Update(permission *model.Permission, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
