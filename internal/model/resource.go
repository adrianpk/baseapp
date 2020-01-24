package model

import (
	"database/sql"
	"strings"

	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
)

type (
	// Resource model
	Resource struct {
		kbs.Identification
		Name        sql.NullString `db:"name" json:"name"`
		Description sql.NullString `db:"description" json:"description"`
		Tag         sql.NullString `db:"tag" json:"tag"`
		Path        sql.NullString `db:"path" json:"path"`
		IsActive    sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted   sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	ResourceForm struct {
		Slug        string `json:"slug" schema:"slug"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		Tag         string `json:"tag" schema:"tag"`
		Path        string `json:"path" schema:"path"`
		IsNew       bool   `json:"-" schema:"-"`
	}
)

func ToResourceFormList(resources []Resource) (fs []ResourceForm) {
	for _, m := range resources {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (resource *Resource) SetCreateValues() error {
	pfx := resource.Name.String
	resource.Identification.SetCreateValues(pfx)
	resource.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (resource *Resource) SetUpdateValues() error {
	resource.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (resource *Resource) Match(tc *Resource) bool {
	r := resource.Identification.Match(tc.Identification) &&
		resource.Name == tc.Name &&
		resource.Description == tc.Description &&
		resource.Tag == tc.Tag &&
		resource.Path == tc.Path
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (resource *Resource) ToForm() ResourceForm {
	return ResourceForm{
		Slug:        resource.Slug.String,
		Name:        resource.Name.String,
		Description: resource.Description.String,
		Tag:         resource.Tag.String,
		Path:        resource.Path.String,
	}
}

// ToModel lets covert a form type to its associated model.
func (resourceForm *ResourceForm) ToModel() Resource {
	tag := strings.ToUpper(resourceForm.Tag)

	return Resource{
		Identification: kbs.Identification{
			Slug: db.ToNullString(resourceForm.Slug),
		},
		Name:        db.ToNullString(resourceForm.Name),
		Description: db.ToNullString(resourceForm.Description),
		Tag:         db.ToNullString(tag),
		Path:        db.ToNullString(resourceForm.Path),
	}
}

func (resourceForm *ResourceForm) GetSlug() string {
	return resourceForm.Slug
}

func (resource *Resource) GenTagIfEmpty() {
	if strings.Trim(resource.Tag.String, " ") == "" {
		resource.Tag = db.ToNullString(kbs.GenTag())
		return
	}
}
