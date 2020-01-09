package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	ResourcePermissionRepo interface {
		Create(u *model.ResourcePermission, tx ...*sqlx.Tx) error
		GetAll() (resourcePermissions []model.ResourcePermission, err error)
		Get(id uuid.UUID) (resourcePermission model.ResourcePermission, err error)
		GetBySlug(slug string) (resourcePermission model.ResourcePermission, err error)
		GetByResourceID(uuid.UUID) ([]model.ResourcePermission, error)
		GetByPermissionID(uuid.UUID) ([]model.ResourcePermission, error)
		Update(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
