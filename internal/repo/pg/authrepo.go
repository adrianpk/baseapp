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
	AuthRepo struct {
		Cfg  *kbs.Config
		Log  kbs.Logger
		Name string
		DB   *sqlx.DB
	}
)

func NewAuthRepo(cfg *kbs.Config, log kbs.Logger, name string, db *sqlx.DB) *AuthRepo {
	return &AuthRepo{
		Cfg:  cfg,
		Log:  log,
		Name: name,
		DB:   db,
	}
}

// ResourcePermission --------------------------------------------------------------------------------
// Create a ResourcePermission
func (ar *AuthRepo) CreateResourcePermission(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error {
	st := `INSERT INTO account_role (id, slug, resource_id, permission_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :resource_id, :permission_id, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	resourcePermission.SetCreateValues()

	_, err = t.NamedExec(st, resourcePermission)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAllResourcePermission
func (ar *AuthRepo) GetAllResourcePermissions() (resourcePermission []model.ResourcePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE is_deleted IS NULL OR NOT is_deleted`

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetResourcePermission account by ID.
func (ar *AuthRepo) GetResourcePermission(id uuid.UUID) (resourcePermission model.ResourcePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetResourcePermissionBySlug account from repo by slug.
func (ar *AuthRepo) GetResourcePermissionBySlug(slug string) (resourcePermission model.ResourcePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&resourcePermission, st)

	return resourcePermission, err
}

// GetResourcePermissiontByResourceID account from repo by slug.
func (ar *AuthRepo) GetResourcePermissionByResourceID(id uuid.UUID) (resourcePermission []model.ResourcePermission, err error) {
	st := `SELECT * FROM accounts WHERE resource_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetResourcePermissiontByRoleID account from repo by slug.
func (ar *AuthRepo) GetResourcePermissionByPermissionID(id uuid.UUID) (resourcePermission []model.ResourcePermission, err error) {
	st := `SELECT * FROM accounts WHERE permission_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// UpdateResourcePermission account data in
func (ar *AuthRepo) UpdateResourcePermission(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error {
	ref, err := ar.GetResourcePermission(resourcePermission.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	resourcePermission.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE resource_permissions SET ")

	if resourcePermission.ResourceID != ref.ResourceID {
		st.WriteString(kbs.SQLStrUpd("resource_id", "resource_id"))
		pcu = true
	}

	if resourcePermission.PermissionID != ref.PermissionID {
		st.WriteString(kbs.SQLStrUpd("permission_id", "permission_id"))
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
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.NamedExec(st.String(), resourcePermission)

	if local {
		ar.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// DeleteResourcePermission account from repo by ID.
func (ar *AuthRepo) DeleteResourcePermission(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM accounts WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// DeleteBySlug account from repo by slug.
func (ar *AuthRepo) DeleteResourcePermissionBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM accounts WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
	st = fmt.Sprintf(st, slug)

	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// RolePermission --------------------------------------------------------------------------------
// Create a RolePermission
func (ar *AuthRepo) CreateRolePermission(resourcePermission *model.RolePermission, tx ...*sqlx.Tx) error {
	st := `INSERT INTO account_role (id, slug, role_id, permission_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :role_id, :permission_id, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	resourcePermission.SetCreateValues()

	_, err = t.NamedExec(st, resourcePermission)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAllRolePermission
func (ar *AuthRepo) GetAllRolePermissions() (resourcePermission []model.RolePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE is_deleted IS NULL OR NOT is_deleted`

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetRolePermission account by ID.
func (ar *AuthRepo) GetRolePermission(id uuid.UUID) (resourcePermission model.RolePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetRolePermissionBySlug account from repo by slug.
func (ar *AuthRepo) GetRolePermissionBySlug(slug string) (resourcePermission model.RolePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&resourcePermission, st)

	return resourcePermission, err
}

// GetRolePermissiontByRoleID account from repo by slug.
func (ar *AuthRepo) GetRolePermissionByRoleID(id uuid.UUID) (resourcePermission []model.RolePermission, err error) {
	st := `SELECT * FROM accounts WHERE role_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetRolePermissiontByRoleID account from repo by slug.
func (ar *AuthRepo) GetRolePermissionByPermissionID(id uuid.UUID) (resourcePermission []model.RolePermission, err error) {
	st := `SELECT * FROM accounts WHERE permission_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// UpdateRolePermission account data in
func (ar *AuthRepo) UpdateRolePermission(resourcePermission *model.RolePermission, tx ...*sqlx.Tx) error {
	ref, err := ar.GetRolePermission(resourcePermission.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	resourcePermission.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE resource_permissions SET ")

	if resourcePermission.RoleID != ref.RoleID {
		st.WriteString(kbs.SQLStrUpd("role_id", "role_id"))
		pcu = true
	}

	if resourcePermission.PermissionID != ref.PermissionID {
		st.WriteString(kbs.SQLStrUpd("permission_id", "permission_id"))
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
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.NamedExec(st.String(), resourcePermission)

	if local {
		ar.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// DeleteRolePermission account from repo by ID.
func (ar *AuthRepo) DeleteRolePermission(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM accounts WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// DeleteBySlug account from repo by slug.
func (ar *AuthRepo) DeleteRolePermissionBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM accounts WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
	st = fmt.Sprintf(st, slug)

	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// AccountRole --------------------------------------------------------------------------------
// Create an AccountRole
func (ar *AuthRepo) CreateAccountRole(accountRole *model.AccountRole, tx ...*sqlx.Tx) error {
	st := `INSERT INTO account_role (id, slug, account_id, role_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :account_id, :role_id, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	accountRole.SetCreateValues()

	_, err = t.NamedExec(st, accountRole)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAllAccountRole
func (ar *AuthRepo) GetAllAccountRoles() (accountRole []model.AccountRole, err error) {
	st := `SELECT * FROM account_roles WHERE is_deleted IS NULL OR NOT is_deleted`

	err = ar.DB.Select(&accountRole, st)
	if err != nil {
		return accountRole, err
	}

	return accountRole, err
}

// GetAccountRole account by ID.
func (ar *AuthRepo) GetAccountRole(id uuid.UUID) (accountRole model.AccountRole, err error) {
	st := `SELECT * FROM account_roles WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&accountRole, st)
	if err != nil {
		return accountRole, err
	}

	return accountRole, err
}

// GetAccountRoleBySlug account from repo by slug.
func (ar *AuthRepo) GetAccountRoleBySlug(slug string) (accountRole model.AccountRole, err error) {
	st := `SELECT * FROM account_roles WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&accountRole, st)

	return accountRole, err
}

// GetAccountRoletByAccountID account from repo by slug.
func (ar *AuthRepo) GetAccountRoleByAccountID(id uuid.UUID) (accountRole []model.AccountRole, err error) {
	st := `SELECT * FROM accounts WHERE account_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&accountRole, st)
	if err != nil {
		return accountRole, err
	}

	return accountRole, err
}

// GetAccountRoletByRoleID account from repo by slug.
func (ar *AuthRepo) GetAccountRoleByRoleID(id uuid.UUID) (accountRole []model.AccountRole, err error) {
	st := `SELECT * FROM accounts WHERE role_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&accountRole, st)
	if err != nil {
		return accountRole, err
	}

	return accountRole, err
}

// UpdateAccountRole account data in
func (ar *AuthRepo) UpdateAccountRole(accountRole *model.AccountRole, tx ...*sqlx.Tx) error {
	ref, err := ar.GetAccountRole(accountRole.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	accountRole.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE account_roles SET ")

	if accountRole.AccountID != ref.AccountID {
		st.WriteString(kbs.SQLStrUpd("account_id", "account_id"))
		pcu = true
	}

	if accountRole.RoleID != ref.RoleID {
		st.WriteString(kbs.SQLStrUpd("role_id", "role_id"))
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
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.NamedExec(st.String(), accountRole)

	if local {
		ar.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// DeleteAccountRole account from repo by ID.
func (ar *AuthRepo) DeleteAccountRole(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM accounts WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	t, local, err := ar.getTx(tx)
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
func (ar *AuthRepo) DeleteAccountRoleBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM accounts WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
	st = fmt.Sprintf(st, slug)

	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	_, err = t.Exec(st)

	if local {
		return t.Commit()
	}

	return err
}

// Misc --------------------------------------------------------------------------------
func (ar *AuthRepo) getTx(txs []*sqlx.Tx) (tx *sqlx.Tx, local bool, err error) {
	// Create a new transaction if its no passed as argument.
	if len(txs) > 0 {
		return txs[0], false, nil
	}

	tx, err = ar.DB.Beginx()
	if err != nil {
		return tx, true, err
	}

	return tx, true, err
}
