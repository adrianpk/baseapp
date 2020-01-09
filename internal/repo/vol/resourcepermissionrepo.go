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
	ResourcePermissionRepo struct {
		Cfg          *kbs.Config
		Log          kbs.Logger
		Name         string
		ResourceID   uuid.UUID
		PermissionID uuid.UUID
	}

	resourcePermissionRow struct {
		mutable bool
		model   model.ResourcePermission
	}
)

var (
	resourcePermission1 = model.ResourcePermission{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("9e7f5355-a8a5-46b7-a3b8-4ddc26c9386b"),
			Slug: db.ToNullString("resourcePermission1-bbc4116229c6"),
		},
		ResourceID:   kbs.ToUUID("e8b43223-17fe-4e36-bd0f-a7d96e867d95"), // userRes
		PermissionID: kbs.ToUUID("00ee4774-776b-4e62-95b1-d32fd248f867"), // permission1
	}

	resourcePermission2 = model.ResourcePermission{
		Identification: kbs.Identification{
			ID:   kbs.ToUUID("de90dce3-1c33-4d79-9dfa-a06fbb7d7c00"),
			Slug: db.ToNullString("resourcePermission2-fd3e9d6aa641"),
		},
		ResourceID:   kbs.ToUUID("fc86c00c-2d4f-400b-ae57-d9d5c87d13c8"), // accountRes
		PermissionID: kbs.ToUUID("00ee4774-776b-4e62-95b1-d32fd248f867"), // permission2
	}

	resourcePermissionsTable = map[uuid.UUID]resourcePermissionRow{
		resourcePermission1.ID: resourcePermissionRow{mutable: false, model: resourcePermission1},
		resourcePermission2.ID: resourcePermissionRow{mutable: false, model: resourcePermission2},
	}
)

func NewResourcePermissionRepo(cfg *kbs.Config, log kbs.Logger, name string) *ResourcePermissionRepo {
	return &ResourcePermissionRepo{
		Cfg:  cfg,
		Log:  log,
		Name: name,
	}
}

// Create a resourcePermission
func (ur *ResourcePermissionRepo) Create(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error {
	_, ok := resourcePermissionsTable[resourcePermission.ID]
	if ok {
		errors.New("duplicate key violates unique constraint")
	}

	if resourcePermission.ID == uuid.Nil {
		errors.New("Non valid primary key")
	}

	resourcePermissionsTable[resourcePermission.ID] = resourcePermissionRow{
		mutable: true,
		model:   *resourcePermission,
	}

	return nil
}

// GetAll resourcePermissionsTable from
func (ur *ResourcePermissionRepo) GetAll() (resourcePermissions []model.ResourcePermission, err error) {
	size := len(resourcePermissionsTable)
	out := make([]model.ResourcePermission, size)
	for _, row := range resourcePermissionsTable {
		out = append(out, row.model)
	}
	return out, nil
}

// Get resourcePermission by ID.
func (ur *ResourcePermissionRepo) Get(id uuid.UUID) (resourcePermission model.ResourcePermission, err error) {
	for _, row := range resourcePermissionsTable {
		if id == row.model.ID {
			return row.model, nil
		}
	}
	return model.ResourcePermission{}, nil
}

// GetBySlug resourcePermission from repo by slug.
func (ur *ResourcePermissionRepo) GetBySlug(slug string) (resourcePermission model.ResourcePermission, err error) {
	for _, row := range resourcePermissionsTable {
		if slug == row.model.Slug.String {
			return row.model, nil
		}
	}
	return model.ResourcePermission{}, nil
}

// GetByResourceID
func (ur *ResourcePermissionRepo) GetByResourceID(resourceID uuid.UUID) (resourcePermissions []model.ResourcePermission, err error) {
	size := len(resourcePermissionsTable)
	resourcePermissions = make([]model.ResourcePermission, size)
	for _, row := range resourcePermissionsTable {
		if resourceID == row.model.ResourceID {
			resourcePermissions = append(resourcePermissions, row.model)
		}
	}
	return resourcePermissions, nil
}

// GetByPermissionID
func (ur *ResourcePermissionRepo) GetByPermissionID(permissionID uuid.UUID) (resourcePermissions []model.ResourcePermission, err error) {
	size := len(resourcePermissionsTable)
	resourcePermissions = make([]model.ResourcePermission, size)
	for _, row := range resourcePermissionsTable {
		if permissionID == row.model.PermissionID {
			resourcePermissions = append(resourcePermissions, row.model)
		}
	}
	return resourcePermissions, nil
}

// Update resourcePermission data in
func (ur *ResourcePermissionRepo) Update(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error {
	for _, row := range resourcePermissionsTable {
		if resourcePermission.ID == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			resourcePermissionsTable[resourcePermission.ID] = resourcePermissionRow{
				mutable: true,
				model:   *resourcePermission,
			}
			return nil
		}
	}
	return errors.New("no records updated")
}

// Delete resourcePermission from repo by ID.
func (ur *ResourcePermissionRepo) Delete(id uuid.UUID, tx ...*sqlx.Tx) error {
	for _, row := range resourcePermissionsTable {
		if id == row.model.ID {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(resourcePermissionsTable, id)
			return nil
		}
	}
	return errors.New("no records deleted")
}

// DeleteBySlug resourcePermission from repo by slug.
func (ur *ResourcePermissionRepo) DeleteBySlug(slug string, tx ...*sqlx.Tx) error {
	for _, row := range resourcePermissionsTable {
		if slug == row.model.Slug.String {
			if !row.mutable {
				return errors.New("non mutable row")
			}

			delete(resourcePermissionsTable, row.model.ID)
			return nil
		}
	}
	return errors.New("no records deleted")
}
