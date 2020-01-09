package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	ResourceRepo interface {
		Create(u *model.Resource, tx ...*sqlx.Tx) error
		GetAll() (resources []model.Resource, err error)
		Get(id uuid.UUID) (resource model.Resource, err error)
		GetBySlug(slug string) (resource model.Resource, err error)
		GetByName(name string) (model.Resource, error)
		GetByTag(tag string) (resource model.Resource, err error)
		GetByPath(path string) (resource model.Resource, err error)
		Update(resource *model.Resource, tx ...*sqlx.Tx) error
		Delete(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
