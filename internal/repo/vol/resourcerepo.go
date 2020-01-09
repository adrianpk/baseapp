package pg

import (
	"errors"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	ResourceRepo struct {
		Cfg  *kbs.Config
		Log  kbs.Logger
		Name string
	}

	resourceRow struct {
		mutable bool
		model   model.Resource
	}
)

var (
	userRes = model.Resource{
		Identification: kbs.Identification{
			ID: kbs.ToUUID("e8b43223-17fe-4e36-bd0f-a7d96e867d95"),
		},
		Name: db.ToNullString("user"),
		Tag:  db.ToNullString("b47f09"),
		Path: db.ToNullString("/users"),
	}

	accountRes = model.Resource{
		Identification: kbs.Identification{
			ID: kbs.ToUUID("fc86c00c-2d4f-400b-ae57-d9d5c87d13c8"),
		},
		Name: db.ToNullString("user"),
		Tag:  db.ToNullString("f0929c"),
		Path: db.ToNullString("/accounts"),
	}

	resourcesTable = map[uuid.UUID]resourceRow{
		userRes.ID:    resourceRow{mutable: false, model: userRes},
		accountRes.ID: resourceRow{mutable: false, model: accountRes},
	}
)

func NewResourceRepo(cfg *kbs.Config, log kbs.Logger, name string) *ResourceRepo {
	return &ResourceRepo{
		Cfg:  cfg,
		Log:  log,
		Name: name,
	}
}

// Create a resource
func (ur *ResourceRepo) Create(resource *model.Resource, tx ...*sqlx.Tx) error {
	_, ok := resourcesTable[resource.ID]
	if ok {
		errors.New("duplicate key violates unique constraint")
	}

	if resource.ID == uuid.Nil {
		errors.New("Non valid primary key")
	}

	resourcesTable[resource.ID] = resourceRow{
		mutable: false,
		model:   *resource,
	}

	return nil
}

// GetAll resourcesTable from
func (ur *ResourceRepo) GetAll() (resources []model.Resource, err error) {
	size := len(resourcesTable)
	out := make([]model.Resource, size)
	for _, row := range resourcesTable {
		out = append(out, row.model)
	}
	return out, nil
}

// Get resource by ID.
func (ur *ResourceRepo) Get(id uuid.UUID) (resource model.Resource, err error) {
	for _, row := range resourcesTable {
		if id == row.model.ID {
			return row.model, nil
		}
	}
	return model.Resource{}, nil
}

// GetBySlug resource from repo by slug.
func (ur *ResourceRepo) GetBySlug(slug string) (resource model.Resource, err error) {
	for _, row := range resourcesTable {
		if slug == row.model.Slug.String {
			return row.model, nil
		}
	}
	return model.Resource{}, nil
}

// GetByName resource from repo by name.
func (ur *ResourceRepo) GetByName(name string) (model.Resource, error) {
	for _, row := range resourcesTable {
		if name == row.model.Name.String {
			return row.model, nil
		}
	}
	return model.Resource{}, nil
}

// GetByTag resource from repo by tag.
func (ur *ResourceRepo) GetByTag(tag string) (model.Resource, error) {
	for _, row := range resourcesTable {
		if tag == row.model.Tag.String {
			return row.model, nil
		}
	}
	return model.Resource{}, nil
}

// Update resource data in
func (ur *ResourceRepo) Update(resource *model.Resource, tx ...*sqlx.Tx) error {
	for _, row := range resourcesTable {
		if resource.ID == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			resourcesTable[resource.ID] = resourceRow{
				mutable: true,
				model:   *resource,
			}
			return nil
		}
	}
	return errors.New("no records updated")
}

// Delete resource from repo by ID.
func (ur *ResourceRepo) Delete(id uuid.UUID, tx ...*sqlx.Tx) error {
	for _, row := range resourcesTable {
		if id == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(resourcesTable, id)
			return nil
		}
	}
	return errors.New("no records deleted")
}

// DeleteBySlug resource from repo by slug.
func (ur *ResourceRepo) DeleteBySlug(slug string, tx ...*sqlx.Tx) error {
	for _, row := range resourcesTable {
		if slug == row.model.Slug.String {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(resourcesTable, row.model.ID)
			return nil
		}
	}
	return errors.New("no records deleted")
}
