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
		Name      sql.NullString `db:"name" json:"name"`
		AccountID uuid.UUID      `db:"account_id" json:"accountID"`
		RoleID    uuid.UUID      `db:"role_id" json:"roleID"`
		IsActive  sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	AccountRoleForm struct {
		Slug      string `json:"slug" schema:"slug"`
		Name      string `json:"name" schema:"name"`
		AccountID string `json:"accountID" schema:"account-id"`
		RoleID    string `json:"roleID" schema:"role-id"`
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
	pfx := accountRole.Name.String
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
		accountRole.Name == tc.Name &&
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
		Name:      accountRole.Name.String,
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
		Name:      db.ToNullString(accountRoleForm.Name),
		AccountID: kbs.ToUUID(accountRoleForm.AccountID),
		RoleID:    kbs.ToUUID(accountRoleForm.RoleID),
	}
}

func (accountRoleForm *AccountRoleForm) GetSlug() string {
	return accountRoleForm.Slug
}
