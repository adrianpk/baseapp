package model

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
)

type (
	// RolePermission model
	RolePermission struct {
		kbs.Identification
		RoleID       uuid.UUID    `db:"resource_id" json:"resourceID"`
		PermissionID uuid.UUID    `db:"permission_id" json:"permissionID"`
		IsActive     sql.NullBool `db:"is_active" json:"isActive"`
		IsDeleted    sql.NullBool `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	RolePermissionForm struct {
		Slug         string `json:"slug" schema:"slug"`
		RoleID       string `json:"resourceID" schema:"resource-id"`
		PermissionID string `json:"permissionID" schema:"permission-id"`
	}
)

func ToRolePermissionFormList(rolePermissions []RolePermission) (fs []RolePermissionForm) {
	for _, m := range rolePermissions {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (rolePermission *RolePermission) SetCreateValues() error {
	pfx := "resourcepermission"
	rolePermission.Identification.SetCreateValues(pfx)
	rolePermission.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (rolePermission *RolePermission) SetUpdateValues() error {
	rolePermission.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (rolePermission *RolePermission) Match(tc *RolePermission) bool {
	r := rolePermission.Identification.Match(tc.Identification) &&
		rolePermission.RoleID == tc.RoleID &&
		rolePermission.PermissionID == tc.PermissionID
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (rolePermission *RolePermission) ToForm() RolePermissionForm {
	return RolePermissionForm{
		RoleID:       rolePermission.RoleID.String(),
		PermissionID: rolePermission.PermissionID.String(),
	}
}

// ToModel lets covert a form type to its associated model.
func (rolePermissionForm *RolePermissionForm) ToModel() RolePermission {
	return RolePermission{
		Identification: kbs.Identification{
			Slug: db.ToNullString(rolePermissionForm.Slug),
		},
		RoleID:       kbs.ToUUID(rolePermissionForm.RoleID),
		PermissionID: kbs.ToUUID(rolePermissionForm.PermissionID),
	}
}

func (rolePermissionForm *RolePermissionForm) GetSlug() string {
	return rolePermissionForm.Slug
}
