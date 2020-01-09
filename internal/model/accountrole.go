package model

import (
	"database/sql"

	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
)

type (
	// AccountRole model
	AccountRole struct {
		kbs.Identification
		AccountID uuid.UUID    `db:"resource_id" json:"resourceID"`
		RoleID    uuid.UUID    `db:"permission_id" json:"permissionID"`
		IsActive  sql.NullBool `db:"is_active" json:"isActive"`
		IsDeleted sql.NullBool `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	AccountRoleForm struct {
		Slug      string `json:"slug" schema:"slug"`
		AccountID string `json:"resourceID" schema:"resource-id"`
		RoleID    string `json:"permissionID" schema:"permission-id"`
	}
)

func ToAccountRoleFormList(accountRoles []AccountRole) (fs []AccountRoleForm) {
	for _, m := range accountRoles {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (accountRole *AccountRole) SetCreateValues() error {
	pfx := "resourcepermission"
	accountRole.Identification.SetCreateValues(pfx)
	accountRole.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (accountRole *AccountRole) SetUpdateValues() error {
	accountRole.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (accountRole *AccountRole) Match(tc *AccountRole) bool {
	r := accountRole.Identification.Match(tc.Identification) &&
		accountRole.AccountID == tc.AccountID &&
		accountRole.RoleID == tc.RoleID
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (accountRole *AccountRole) ToForm() AccountRoleForm {
	return AccountRoleForm{
		AccountID: accountRole.AccountID.String(),
		RoleID:    accountRole.RoleID.String(),
	}
}

// ToModel lets covert a form type to its associated model.
func (accountRoleForm *AccountRoleForm) ToModel() AccountRole {
	return AccountRole{
		Identification: kbs.Identification{
			Slug: db.ToNullString(accountRoleForm.Slug),
		},
		AccountID: kbs.ToUUID(accountRoleForm.AccountID),
		RoleID:    kbs.ToUUID(accountRoleForm.RoleID),
	}
}

func (accountRoleForm *AccountRoleForm) GetSlug() string {
	return accountRoleForm.Slug
}
