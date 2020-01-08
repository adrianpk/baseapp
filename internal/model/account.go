package model

import (
	"database/sql"

	"github.com/lib/pq"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"

	uuid "github.com/satori/go.uuid"
)

type (
	// User model
	Account struct {
		kbs.Identification
		OwnerID     uuid.UUID      `db:"owner_id" json:"ownerID"`
		ParentID    uuid.UUID      `db:"parent_id" json:"parentID"`
		AccountType sql.NullString `db:"account_type" json:"accountType"`
		Name        sql.NullString `db:"name" json:"name"`
		Email       sql.NullString `db:"email" json:"email"`
		Locale      sql.NullString `db:"locale" json:"locale"`
		BaseTZ      sql.NullString `db:"base_tz" json:"baseTZ"`
		CurrentTZ   sql.NullString `db:"current_tz" json:"currentTZ"`
		StartsAt    pq.NullTime    `db:"starts_at" json:"startsAt"`
		EndsAt      pq.NullTime    `db:"ends_at" json:"endsAt"`
		IsActive    sql.NullBool   `db:"is_active" json:"isActive"`
		IsDeleted   sql.NullBool   `db:"is_deleted" json:"isDeleted"`
		kbs.Audit
	}
)

type (
	AccountForm struct {
		Slug        string `json:"slug" schema:"slug"`
		OwnerID     string `json:"ownerID" schema:"owner-id"`
		ParentID    string `json:"parentID" schema:"parent-id"`
		AccountType string `json:"accountType" schema:"account-type"`
		Name        string `json:"name" schema:"name"`
		Email       string `json:"email" schema:"email"`
		IsNew       bool   `json:"-" schema:"-"`
	}
)

func ToAccountFormList(accounts []Account) (fs []AccountForm) {
	for _, m := range accounts {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// SetCreateValues sets de ID, slug and audit values.
func (account *Account) SetCreateValues() error {
	pfx := account.Name.String
	account.Identification.SetCreateValues(pfx)
	account.Audit.SetCreateValues()
	return nil
}

// SetUpdateValues updates audit values.
func (account *Account) SetUpdateValues() error {
	account.Audit.SetUpdateValues()
	return nil
}

// Match condition for
func (account *Account) Match(tc *Account) bool {
	r := account.Identification.Match(tc.Identification) &&
		account.OwnerID == tc.OwnerID &&
		account.ParentID == tc.ParentID &&
		account.AccountType == tc.AccountType &&
		account.Name == tc.Name &&
		account.Email == tc.Email
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (account *Account) ToForm() AccountForm {
	return AccountForm{
		OwnerID:     account.OwnerID.String(),
		ParentID:    account.ParentID.String(),
		AccountType: account.AccountType.String,
		Name:        account.Name.String,
		Email:       account.Email.String,
	}
}

// ToModel lets covert a form type to its associated model.
func (accountForm *AccountForm) ToModel() Account {
	return Account{
		Identification: kbs.Identification{
			Slug: db.ToNullString(accountForm.Slug),
		},
		OwnerID:     kbs.ToUUID(accountForm.OwnerID),
		ParentID:    kbs.ToUUID(accountForm.ParentID),
		AccountType: db.ToNullString(accountForm.AccountType),
		Name:        db.ToNullString(accountForm.Name),
		Email:       db.ToNullString(accountForm.Email),
	}
}

func (accountForm *AccountForm) GetSlug() string {
	return accountForm.Slug
}
