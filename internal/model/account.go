package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"

	uuid "github.com/satori/go.uuid"
)

type (
	// Account model
	Account struct {
		kbs.Identification
		OwnerID     uuid.UUID      `db:"owner_id" json:"ownerID"`
		ParentID    uuid.UUID      `db:"parent_id" json:"parentID"`
		AccountType sql.NullString `db:"account_type" json:"accountType"`
		Username    sql.NullString `db:"username" json:"username"`
		Email       sql.NullString `db:"email" json:"email"`
		GivenName   sql.NullString `db:"given_name" json:"givenName" schema:"given-name"`
		MiddleNames sql.NullString `db:"middle_names" json:"middleNames" schema:"middle-names"`
		FamilyName  sql.NullString `db:"family_name" json:"familyName" schema:"family-name"`
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
		Username    string `json:"username" schema:"username"`
		Email       string `json:"email" schema:"email"`
		GivenName   string `json:"givenName" schema:"given-name"`
		MiddleNames string `json:"middleNames" schema:"middle-names"`
		FamilyName  string `json:"familyName" schema:"family-name"`
		IsNew       bool   `json:"-" schema:"-"`
	}
)

func ToAccountFormList(accounts []Account) (fs []AccountForm) {
	for _, m := range accounts {
		fs = append(fs, m.ToForm())
	}
	return fs
}

func (account *Account) FullName() string {
	return strings.Trim(fmt.Sprintf("%s %s", account.GivenName.String, account.FamilyName.String), " ")
}

// SetCreateValues sets de ID, slug and audit values.
func (account *Account) SetCreateValues() error {
	pfx := account.Username.String
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
		account.Username == tc.Username &&
		account.Email == tc.Email &&
		account.GivenName == tc.GivenName &&
		account.MiddleNames == tc.MiddleNames &&
		account.FamilyName == tc.FamilyName
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
		Username:    account.Username.String,
		Email:       account.Email.String,
		GivenName:   account.GivenName.String,
		MiddleNames: account.MiddleNames.String,
		FamilyName:  account.FamilyName.String,
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
		Username:    db.ToNullString(accountForm.Username),
		Email:       db.ToNullString(accountForm.Email),
		GivenName:   db.ToNullString(accountForm.GivenName),
		MiddleNames: db.ToNullString(accountForm.MiddleNames),
		FamilyName:  db.ToNullString(accountForm.FamilyName),
	}
}

func (accountForm AccountForm) GetSlug() string {
	return accountForm.Slug
}
