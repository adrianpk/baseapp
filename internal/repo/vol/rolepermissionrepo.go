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
	rolePermissionRow struct {
		mutable bool
		model   model.RolePermission
	}
)

var (
	rolePermission1 = model.RolePermission{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("9e7f5355-a8a5-46b7-a3b8-4ddc26c9386b"),
			Slug: db.ToNullString("rolePermission1-bbc4116229c6"),
		},
		RoleID:       kbs.ToUUID("288bb973-2196-4007-808a-d7844ecf4dd9"), // userRes
		PermissionID: kbs.ToUUID("00ee4774-776b-4e62-95b1-d32fd248f867"), // permission1
	}

	rolePermission2 = model.RolePermission{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("de90dce3-1c33-4d79-9dfa-a06fbb7d7c00"),
			Slug: db.ToNullString("rolePermission2-fd3e9d6aa641"),
		},
		RoleID:       kbs.ToUUID("d0d6bc3a-38b0-4a00-83c0-516d2514d7b5"), // accountRes
		PermissionID: kbs.ToUUID("2c9bba14-c579-4c44-a2da-6ff15324605c"), // permission2
	}

	rolePermissionsTable = map[uuid.UUID]rolePermissionRow{
		rolePermission1.ID: rolePermissionRow{mutable: false, model: rolePermission1},
		rolePermission2.ID: rolePermissionRow{mutable: false, model: rolePermission2},
	}
)

func (ar *AuthRepo) CreateRolePermission(rolePermission *model.RolePermission, tx ...*sqlx.Tx) error {
	_, ok := rolePermissionsTable[rolePermission.ID]
	if ok {
		errors.New("duplicate key violates unique constraint")
	}

	if rolePermission.ID == uuid.Nil {
		errors.New("Non valid primary key")
	}

	rolePermissionsTable[rolePermission.ID] = rolePermissionRow{
		mutable: true,
		model:   *rolePermission,
	}

	return nil
}

func (ar *AuthRepo) GetAllRolePermissions() (rolePermissions []model.RolePermission, err error) {
	size := len(rolePermissionsTable)
	out := make([]model.RolePermission, size)
	for _, row := range rolePermissionsTable {
		out = append(out, row.model)
	}
	return out, nil
}

func (ar *AuthRepo) GetRolePermission(id uuid.UUID) (rolePermission model.RolePermission, err error) {
	for _, row := range rolePermissionsTable {
		if id == row.model.ID {
			return row.model, nil
		}
	}
	return model.RolePermission{}, nil
}

func (ar *AuthRepo) GetRolePermissionBySlug(slug string) (rolePermission model.RolePermission, err error) {
	for _, row := range rolePermissionsTable {
		if slug == row.model.Slug.String {
			return row.model, nil
		}
	}
	return model.RolePermission{}, nil
}

func (ar *AuthRepo) GetRolePermissionByRoleID(resourceID uuid.UUID) (rolePermissions []model.RolePermission, err error) {
	size := len(rolePermissionsTable)
	rolePermissions = make([]model.RolePermission, size)
	for _, row := range rolePermissionsTable {
		if resourceID == row.model.RoleID {
			rolePermissions = append(rolePermissions, row.model)
		}
	}
	return rolePermissions, nil
}

func (ar *AuthRepo) GetRolePermissionByPermissionID(permissionID uuid.UUID) (rolePermissions []model.RolePermission, err error) {
	size := len(rolePermissionsTable)
	rolePermissions = make([]model.RolePermission, size)
	for _, row := range rolePermissionsTable {
		if permissionID == row.model.PermissionID {
			rolePermissions = append(rolePermissions, row.model)
		}
	}
	return rolePermissions, nil
}

func (ar *AuthRepo) UpdateRolePermission(rolePermission *model.RolePermission, tx ...*sqlx.Tx) error {
	for _, row := range rolePermissionsTable {
		if rolePermission.ID == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			rolePermissionsTable[rolePermission.ID] = rolePermissionRow{
				mutable: true,
				model:   *rolePermission,
			}
			return nil
		}
	}
	return errors.New("no records updated")
}

func (ar *AuthRepo) DeleteRolePermission(id uuid.UUID, tx ...*sqlx.Tx) error {
	for _, row := range rolePermissionsTable {
		if id == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(rolePermissionsTable, id)
			return nil
		}
	}
	return errors.New("no records deleted")
}

func (ar *AuthRepo) DeleteRolePermissionBySlug(slug string, tx ...*sqlx.Tx) error {
	for _, row := range rolePermissionsTable {
		if slug == row.model.Slug.String {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(rolePermissionsTable, row.model.ID)
			return nil
		}
	}
	return errors.New("no records deleted")
}

// Custom

func (ar *AuthRepo) GetRolePermissions(slug string) (permissions []model.Permission, err error) {
	panic("not implemented")
}

func (ar *AuthRepo) GetNotRolePermissions(slug string) (permissions []model.Permission, err error) {
	panic("not implemented")
}

func (ar *AuthRepo) AppendRolePermission(roleSlug, permissionSlug string) (err error) {
	panic("not implemented")
}

func (ar *AuthRepo) RemoveRolePermission(roleSlug, permissionSlug string) (err error) {
	panic("not implemented")
}
