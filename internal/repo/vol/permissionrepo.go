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
	permissionRow struct {
		mutable bool
		model   model.Permission
	}
)

var (
	permission1 = model.Permission{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("00ee4774-776b-4e62-95b1-d32fd248f867"),
			Slug: db.ToNullString("permission1-c4c55224c7d1"),
		},
		Name:        db.ToNullString("permission1"),
		Description: db.ToNullString("permission1 description"),
	}

	permission2 = model.Permission{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("2c9bba14-c579-4c44-a2da-6ff15324605c"),
			Slug: db.ToNullString("permission2-0ac0d549ab01"),
		},
		Name:        db.ToNullString("permission1"),
		Description: db.ToNullString("permission1 description"),
	}

	permissionsTable = map[uuid.UUID]permissionRow{
		permission1.ID: permissionRow{mutable: false, model: permission1},
		permission2.ID: permissionRow{mutable: false, model: permission2},
	}
)

func (ar *AuthRepo) CreatePermission(permission *model.Permission, tx ...*sqlx.Tx) error {
	_, ok := permissionsTable[permission.ID]
	if ok {
		errors.New("duplicate key violates unique constraint")
	}

	if permission.ID == uuid.Nil {
		errors.New("Non valid primary key")
	}

	permissionsTable[permission.ID] = permissionRow{
		mutable: true,
		model:   *permission,
	}

	return nil
}

func (ar *AuthRepo) GetAllPermissions() (permissions []model.Permission, err error) {
	size := len(permissionsTable)
	out := make([]model.Permission, size)
	for _, row := range permissionsTable {
		out = append(out, row.model)
	}
	return out, nil
}

func (ar *AuthRepo) GetPermission(id uuid.UUID) (permission model.Permission, err error) {
	for _, row := range permissionsTable {
		if id == row.model.ID {
			return row.model, nil
		}
	}
	return model.Permission{}, nil
}

func (ar *AuthRepo) GetPermissionBySlug(slug string) (permission model.Permission, err error) {
	for _, row := range permissionsTable {
		if slug == row.model.Slug.String {
			return row.model, nil
		}
	}
	return model.Permission{}, nil
}

func (ar *AuthRepo) GetPermissionByName(name string) (model.Permission, error) {
	for _, row := range permissionsTable {
		if name == row.model.Name.String {
			return row.model, nil
		}
	}
	return model.Permission{}, nil
}

func (ar *AuthRepo) UpdatePermission(permission *model.Permission, tx ...*sqlx.Tx) error {
	for _, row := range permissionsTable {
		if permission.ID == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			permissionsTable[permission.ID] = permissionRow{
				mutable: true,
				model:   *permission,
			}
			return nil
		}
	}
	return errors.New("no records updated")
}

func (ar *AuthRepo) DeletePermission(id uuid.UUID, tx ...*sqlx.Tx) error {
	for _, row := range permissionsTable {
		if id == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(permissionsTable, id)
			return nil
		}
	}
	return errors.New("no records deleted")
}

func (ar *AuthRepo) DeletePermissionBySlug(slug string, tx ...*sqlx.Tx) error {
	for _, row := range permissionsTable {
		if slug == row.model.Slug.String {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(permissionsTable, row.model.ID)
			return nil
		}
	}
	return errors.New("no records deleted")
}
