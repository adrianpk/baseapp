package model

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
)

type (
	// ResourcePermission model
	ResourcePermission struct {
		kbs.Identification
		Name         sql.NullString `db:"name" json:"name"`
		ResourceID   uuid.UUID      `db:"resource_id" json:"resourceID"`
		PermissionID uuid.UUID      `db:"permission_id" json:"permissionID"`
		IsActive     sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted    sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	ResourcePermissionForm struct {
		Slug         string `json:"slug" schema:"slug"`
		Name         string `json:"name" schema:"name"`
		ResourceID   string `json:"resourceID" schema:"resource-id"`
		PermissionID string `json:"permissionID" schema:"permission-id"`
	}
)

func ToResourcePermissionFormList(resourcePermissions []ResourcePermission) (fs []ResourcePermissionForm) {
	for _, m := range resourcePermissions {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (resourcePermission *ResourcePermission) SetCreateValues() error {
	pfx := "resourcepermission"
	resourcePermission.Identification.SetCreateValues(pfx)
	resourcePermission.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (resourcePermission *ResourcePermission) SetUpdateValues() error {
	resourcePermission.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (resourcePermission *ResourcePermission) Match(tc *ResourcePermission) bool {
	r := resourcePermission.Identification.Match(tc.Identification) &&
		resourcePermission.Name == tc.Name &&
		resourcePermission.ResourceID == tc.ResourceID &&
		resourcePermission.PermissionID == tc.PermissionID
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (resourcePermission *ResourcePermission) ToForm() ResourcePermissionForm {
	return ResourcePermissionForm{
		Name:         resourcePermission.Name.String,
		ResourceID:   resourcePermission.ResourceID.String(),
		PermissionID: resourcePermission.PermissionID.String(),
	}
}

// ToModel lets covert a form type to its associated model.
func (resourcePermissionForm *ResourcePermissionForm) ToModel() ResourcePermission {
	return ResourcePermission{
		Identification: kbs.Identification{
			Slug: db.ToNullString(resourcePermissionForm.Slug),
		},
		Name:         db.ToNullString(resourcePermissionForm.Name),
		ResourceID:   kbs.ToUUID(resourcePermissionForm.ResourceID),
		PermissionID: kbs.ToUUID(resourcePermissionForm.PermissionID),
	}
}

func (resourcePermissionForm ResourcePermissionForm) GetSlug() string {
	return resourcePermissionForm.Slug
}
