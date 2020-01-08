package pg

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	AccountRepo struct {
		Cfg  *kbs.Config
		Log  kbs.Logger
		Name string
		DB   *sqlx.DB
	}
)

func NewAccountRepo(cfg *kbs.Config, log kbs.Logger, name string, db *sqlx.DB) *AccountRepo {
	return &AccountRepo{
		Cfg:  cfg,
		Log:  log,
		Name: name,
		DB:   db,
	}
}

// Create a account
func (ur *AccountRepo) Create(account *model.Account, tx ...*sqlx.Tx) error {
	st := `INSERT INTO accounts (id, slug, owner_id, parent_id, account_type, name, email, locale, base_tz, current_tz, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :owner_id, :parent_id, :account_type, :name, :email, :locale, :base_tz, :current_tz, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ur.getTx(tx)
	if err != nil {
		return err
	}

	account.SetCreateValues()

	_, err = t.NamedExec(st, account)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAll accounts from
func (ur *AccountRepo) GetAll() (accounts []model.Account, err error) {
	st := `SELECT * FROM accounts WHERE is_delete IS NULL OR NOT is_deleted`

	err = ur.DB.Select(&accounts, st)
	if err != nil {
		return accounts, err
	}

	return accounts, err
}

// Get account by ID.
func (ur *AccountRepo) Get(id uuid.UUID) (account model.Account, err error) {
	st := `SELECT * FROM USERS WHERE id = '%s' WHERE is_delete IS NULL OR NOT is_deleted LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ur.DB.Get(&account, st)
	if err != nil {
		return account, err
	}

	return account, err
}

// GetBySlug account from repo by slug.
func (ur *AccountRepo) GetBySlug(slug string) (account model.Account, err error) {
	st := `SELECT * FROM USERS WHERE slug = '%s' WHERE is_delete IS NULL OR NOT is_deleted LIMIT 1;`
	st = fmt.Sprintf(st, slug)

	err = ur.DB.Get(&account, st)

	return account, err
}

// GetByName account from repo by accountname.
func (ur *AccountRepo) GetByName(name string) (model.Account, error) {
	var account model.Account

	st := `SELECT * FROM USERS WHERE name = '%s' LIMIT es_deleted IS NULL or NOT is_deleted LIMIT 1;`
	st = fmt.Sprintf(st, name)

	err := ur.DB.Get(&account, st)

	return account, err
}

// Update account data in
func (ur *AccountRepo) Update(account *model.Account, tx ...*sqlx.Tx) error {
	ref, err := ur.Get(account.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	account.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE accounts SET ")

	if account.OwnerID != ref.OwnerID {
		st.WriteString(kbs.SQLStrUpd("owner_id", "owner_id"))
		pcu = true
	}

	if account.ParentID != ref.ParentID {
		st.WriteString(kbs.SQLStrUpd("parent_id", "parent_id"))
		pcu = true
	}

	if account.Name != ref.Name {
		st.WriteString(kbs.SQLComma(pcu))
		st.WriteString(kbs.SQLStrUpd("name", "name"))
		pcu = true
	}

	if account.Email != ref.Email {
		st.WriteString(kbs.SQLComma(pcu))
		st.WriteString(kbs.SQLStrUpd("email", "email"))
		pcu = true
	}

	st.WriteString(" ")
	st.WriteString(kbs.SQLWhereID(ref.ID.String()))
	st.WriteString(";")

	//fmt.Println(st.String())

	if pcu == false {
		return errors.New("no fields to update")
	}

	// Create a local transaction if it is not passed as argument.
	t, local, err := ur.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.NamedExec(st.String(), account)

	if local {
		ur.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// Delete account from repo by ID.
func (ur *AccountRepo) Delete(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM USERS WHERE id = '%s';`
	st = fmt.Sprintf(st, id)

	t, local, err := ur.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// DeleteBySlug:w account from repo by slug.
func (ur *AccountRepo) DeleteBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM USERS WHERE slug = '%s';`
	st = fmt.Sprintf(st, slug)

	t, local, err := ur.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// DeleteByaccountname account from repo by accountname.
func (ur *AccountRepo) DeleteByAccountname(accountname string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM USERS WHERE accountname = '%s';`
	st = fmt.Sprintf(st, accountname)

	t, local, err := ur.getTx(tx)
	if err != nil {
		return err
	}
	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// GetBySlug account from repo by slug token.
func (ur *AccountRepo) GetBySlugAndToken(slug, token string) (model.Account, error) {
	var account model.Account

	st := `SELECT * FROM USERS WHERE slug = '%s' AND confirmation_token = '%s' LIMIT 1;`
	st = fmt.Sprintf(st, slug, token)

	err := ur.DB.Get(&account, st)

	return account, err
}

// Confirm account from repo by slug.
func (ur *AccountRepo) ConfirmAccount(slug, token string, tx ...*sqlx.Tx) (err error) {
	st := `UPDATE USERS SET is_confirmed = TRUE WHERE slug = '%s' AND confirmation_token = '%s';`
	st = fmt.Sprintf(st, slug, token)

	t, local, err := ur.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

func (ur *AccountRepo) newTx() (tx *sqlx.Tx, err error) {
	tx, err = ur.DB.Beginx()
	if err != nil {
		return tx, err
	}

	return tx, err
}

func (ur *AccountRepo) getTx(txs []*sqlx.Tx) (tx *sqlx.Tx, local bool, err error) {
	// Create a new transaction if its no passed as argument.
	if len(txs) > 0 {
		return txs[0], false, nil
	}

	tx, err = ur.DB.Beginx()
	if err != nil {
		return tx, true, err
	}

	return tx, true, err
}
