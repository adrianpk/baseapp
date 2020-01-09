package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	RolePermissionRepo interface {
		Create(u *model.RolePermission, tx ...*sqlx.Tx) error
		GetAll() (rolePermissions []model.RolePermission, err error)
		Get(id uuid.UUID) (rolePermission model.RolePermission, err error)
		GetBySlug(slug string) (rolePermission model.RolePermission, err error)
		GetByRoleID(uuid.UUID) ([]model.RolePermission, error)
		GetByPermissionID(uuid.UUID) ([]model.RolePermission, error)
		Update(rolePermission *model.RolePermission, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
