package model

import (
	"database/sql"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
)

type (
	// Role model
	Role struct {
		kbs.Identification
		Name        sql.NullString `db:"name" json:"name"`
		Description sql.NullString `db:"description" json:"description"`
		IsActive    sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted   sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	RoleForm struct {
		Slug        string `json:"slug" schema:"slug"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		Tag         string `json:"tag" schema:"tag"`
		Path        string `json:"path" schema:"path"`
		IsNew       bool   `json:"-" schema:"-"`
	}
)

func ToRoleFormList(roles []Role) (fs []RoleForm) {
	for _, m := range roles {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (role *Role) SetCreateValues() error {
	pfx := role.Name.String
	role.Identification.SetCreateValues(pfx)
	role.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (role *Role) SetUpdateValues() error {
	role.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (role *Role) Match(tc *Role) bool {
	r := role.Identification.Match(tc.Identification) &&
		role.Name == tc.Name &&
		role.Description == tc.Description
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (role *Role) ToForm() RoleForm {
	return RoleForm{
		Name:        role.Name.String,
		Description: role.Description.String,
	}
}

// ToModel lets covert a form type to its associated model.
func (roleForm *RoleForm) ToModel() Role {
	return Role{
		Identification: kbs.Identification{
			Slug: db.ToNullString(roleForm.Slug),
		},
		Name:        db.ToNullString(roleForm.Name),
		Description: db.ToNullString(roleForm.Description),
	}
}

func (roleForm RoleForm) GetSlug() string {
	return roleForm.Slug
}
