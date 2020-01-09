package model

import (
	"database/sql"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
)

type (
	// Permission model
	Permission struct {
		kbs.Identification
		Name        sql.NullString `db:"name" json:"name"`
		Description sql.NullString `db:"description" json:"description"`
		IsActive    sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted   sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	PermissionForm struct {
		Slug        string `json:"slug" schema:"slug"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		Tag         string `json:"tag" schema:"tag"`
		Path        string `json:"path" schema:"path"`
		IsNew       bool   `json:"-" schema:"-"`
	}
)

func ToPermissionFormList(permissions []Permission) (fs []PermissionForm) {
	for _, m := range permissions {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (permission *Permission) SetCreateValues() error {
	pfx := permission.Name.String
	permission.Identification.SetCreateValues(pfx)
	permission.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (permission *Permission) SetUpdateValues() error {
	permission.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (permission *Permission) Match(tc *Permission) bool {
	r := permission.Identification.Match(tc.Identification) &&
		permission.Name == tc.Name &&
		permission.Description == tc.Description
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (permission *Permission) ToForm() PermissionForm {
	return PermissionForm{
		Name:        permission.Name.String,
		Description: permission.Description.String,
	}
}

// ToModel lets covert a form type to its associated model.
func (permissionForm *PermissionForm) ToModel() Permission {
	return Permission{
		Identification: kbs.Identification{
			Slug: db.ToNullString(permissionForm.Slug),
		},
		Name:        db.ToNullString(permissionForm.Name),
		Description: db.ToNullString(permissionForm.Description),
	}
}

func (permissionForm *PermissionForm) GetSlug() string {
	return permissionForm.Slug
}
