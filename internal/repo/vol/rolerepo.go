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
	RoleRepo struct {
		Cfg  *kbs.Config
		Log  kbs.Logger
		Name string
	}

	roleRow struct {
		mutable bool
		model   model.Role
	}
)

var (
	role1 = model.Role{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("288bb973-2196-4007-808a-d7844ecf4dd9"),
			Slug: db.ToNullString("role1-6ccf99f1a582"),
		},
		Name:        db.ToNullString("role1"),
		Description: db.ToNullString("role1 description"),
	}

	role2 = model.Role{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("d0d6bc3a-38b0-4a00-83c0-516d2514d7b5"),
			Slug: db.ToNullString("role2-2de6909780aa"),
		},
		Name:        db.ToNullString("role1"),
		Description: db.ToNullString("role1 description"),
	}

	rolesTable = map[uuid.UUID]roleRow{
		role1.ID: roleRow{mutable: false, model: role1},
		role2.ID: roleRow{mutable: false, model: role2},
	}
)

func NewRoleRepo(cfg *kbs.Config, log kbs.Logger, name string) *RoleRepo {
	return &RoleRepo{
		Cfg:  cfg,
		Log:  log,
		Name: name,
	}
}

// Create a role
func (ur *RoleRepo) Create(role *model.Role, tx ...*sqlx.Tx) error {
	_, ok := rolesTable[role.ID]
	if ok {
		errors.New("duplicate key violates unique constraint")
	}

	if role.ID == uuid.Nil {
		errors.New("Non valid primary key")
	}

	rolesTable[role.ID] = roleRow{
		mutable: true,
		model:   *role,
	}

	return nil
}

// GetAll rolesTable from
func (ur *RoleRepo) GetAll() (roles []model.Role, err error) {
	size := len(rolesTable)
	out := make([]model.Role, size)
	for _, row := range rolesTable {
		out = append(out, row.model)
	}
	return out, nil
}

// Get role by ID.
func (ur *RoleRepo) Get(id uuid.UUID) (role model.Role, err error) {
	for _, row := range rolesTable {
		if id == row.model.ID {
			return row.model, nil
		}
	}
	return model.Role{}, nil
}

// GetBySlug role from repo by slug.
func (ur *RoleRepo) GetBySlug(slug string) (role model.Role, err error) {
	for _, row := range rolesTable {
		if slug == row.model.Slug.String {
			return row.model, nil
		}
	}
	return model.Role{}, nil
}

// GetByName role from repo by name.
func (ur *RoleRepo) GetByName(name string) (model.Role, error) {
	for _, row := range rolesTable {
		if name == row.model.Name.String {
			return row.model, nil
		}
	}
	return model.Role{}, nil
}

// Update role data in
func (ur *RoleRepo) Update(role *model.Role, tx ...*sqlx.Tx) error {
	for _, row := range rolesTable {
		if role.ID == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			rolesTable[role.ID] = roleRow{
				mutable: true,
				model:   *role,
			}
			return nil
		}
	}
	return errors.New("no records updated")
}

// Delete role from repo by ID.
func (ur *RoleRepo) Delete(id uuid.UUID, tx ...*sqlx.Tx) error {
	for _, row := range rolesTable {
		if id == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(rolesTable, id)
			return nil
		}
	}
	return errors.New("no records deleted")
}

// DeleteBySlug role from repo by slug.
func (ur *RoleRepo) DeleteBySlug(slug string, tx ...*sqlx.Tx) error {
	for _, row := range rolesTable {
		if slug == row.model.Slug.String {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(rolesTable, row.model.ID)
			return nil
		}
	}
	return errors.New("no records deleted")
}
