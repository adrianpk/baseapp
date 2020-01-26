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

// Resource --------------------------------------------------------------------------------
// Create a Resource
func (ar *AuthRepo) CreateResource(resource *model.Resource, tx ...*sqlx.Tx) error {
	st := `INSERT INTO resources (id, slug, tenant_id, name, description, tag, path, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :name, :description, :tag, :path, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at);`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	resource.SetCreateValues()

	_, err = t.NamedExec(st, resource)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAllResource
func (ar *AuthRepo) GetAllResources() (resource []model.Resource, err error) {
	st := `SELECT * FROM resources WHERE is_deleted IS NULL OR NOT is_deleted;`

	err = ar.DB.Select(&resource, st)
	if err != nil {
		return resource, err
	}

	return resource, err
}

// GetResource account by ID.
func (ar *AuthRepo) GetResource(id uuid.UUID) (resource model.Resource, err error) {
	st := `SELECT * FROM resources WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&resource, st)
	if err != nil {
		return resource, err
	}

	return resource, err
}

// GetResourceBySlug account from repo by slug.
func (ar *AuthRepo) GetResourceBySlug(slug string) (resource model.Resource, err error) {
	st := `SELECT * FROM resources WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&resource, st)

	return resource, err
}

// GetResourceByName account from repo by slug.
func (ar *AuthRepo) GetResourceByName(name string) (resource model.Resource, err error) {
	st := `SELECT * FROM resources WHERE name = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, name)

	err = ar.DB.Get(&resource, st)

	return resource, err
}

// GetResourceByTag account from repo by tag.
func (ar *AuthRepo) GetResourceByTag(tag string) (resource model.Resource, err error) {
	st := `SELECT * FROM resources WHERE tag = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, tag)

	err = ar.DB.Get(&resource, st)

	return resource, err
}

// GetResourceByPath account from repo by path.
func (ar *AuthRepo) GetResourceByPath(path string) (resource model.Resource, err error) {
	st := `SELECT * FROM resources WHERE path = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, path)

	err = ar.DB.Get(&resource, st)

	return resource, err
}

// Update resource
func (ar *AuthRepo) UpdateResource(resource *model.Resource, tx ...*sqlx.Tx) error {
	ref, err := ar.GetResource(resource.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	resource.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE resources SET ")

	if resource.Name.String != ref.Name.String {
		st.WriteString(kbs.SQLStrUpd("name", "name"))
		pcu = true
	}

	if resource.Tag.String != ref.Tag.String {
		st.WriteString(kbs.SQLStrUpd("tag", "tag"))
		pcu = true
	}

	if resource.Path.String != ref.Path.String {
		st.WriteString(kbs.SQLStrUpd("path", "path"))
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

	_, err = t.NamedExec(st.String(), resource)

	if local {
		ar.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// DeleteResource account from repo by ID.
func (ar *AuthRepo) DeleteResource(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM resource WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
func (ar *AuthRepo) DeleteResourceBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM resources WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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

// Permission --------------------------------------------------------------------------------
// Create a Permission
func (ar *AuthRepo) CreatePermission(permission *model.Permission, tx ...*sqlx.Tx) error {
	st := `INSERT INTO permission (id, slug, name, description, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :name, :description, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	permission.SetCreateValues()

	_, err = t.NamedExec(st, permission)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAllPermission
func (ar *AuthRepo) GetAllPermissions() (permission []model.Permission, err error) {
	st := `SELECT * FROM permissions WHERE is_deleted IS NULL OR NOT is_deleted`

	err = ar.DB.Select(&permission, st)
	if err != nil {
		return permission, err
	}

	return permission, err
}

// GetPermission account by ID.
func (ar *AuthRepo) GetPermission(id uuid.UUID) (permission model.Permission, err error) {
	st := `SELECT * FROM permissions WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&permission, st)
	if err != nil {
		return permission, err
	}

	return permission, err
}

// GetPermissionBySlug account from repo by slug.
func (ar *AuthRepo) GetPermissionBySlug(slug string) (permission model.Permission, err error) {
	st := `SELECT * FROM permissions WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&permission, st)

	return permission, err
}

// GetPermissionByName account from repo by slug.
func (ar *AuthRepo) GetPermissionByName(name string) (permission model.Permission, err error) {
	st := `SELECT * FROM permissions WHERE name = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, name)

	err = ar.DB.Get(&permission, st)

	return permission, err
}

// Update permission
func (ar *AuthRepo) UpdatePermission(permission *model.Permission, tx ...*sqlx.Tx) error {
	ref, err := ar.GetPermission(permission.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	permission.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE permissions SET ")

	if permission.Name.String != ref.Name.String {
		st.WriteString(kbs.SQLStrUpd("name", "name"))
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

	_, err = t.NamedExec(st.String(), permission)

	if local {
		ar.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// DeletePermission account from repo by ID.
func (ar *AuthRepo) DeletePermission(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM permission WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
func (ar *AuthRepo) DeletePermissionBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM permissions WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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

// Role --------------------------------------------------------------------------------
// Create a Role
func (ar *AuthRepo) CreateRole(role *model.Role, tx ...*sqlx.Tx) error {
	st := `INSERT INTO roles (id, slug, name, description, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :name, :description, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// Create a local transaction if it is not passed as argument.
	t, local, err := ar.getTx(tx)
	if err != nil {
		return err
	}

	role.SetCreateValues()

	_, err = t.NamedExec(st, role)
	if err != nil {
		return err
	}

	// Commit on local transactions
	if local {
		return t.Commit()
	}

	return nil
}

// GetAllRole
func (ar *AuthRepo) GetAllRoles() (role []model.Role, err error) {
	st := `SELECT * FROM roles WHERE is_deleted IS NULL OR NOT is_deleted`

	err = ar.DB.Select(&role, st)
	if err != nil {
		return role, err
	}

	return role, err
}

// GetRole account by ID.
func (ar *AuthRepo) GetRole(id uuid.UUID) (role model.Role, err error) {
	st := `SELECT * FROM roles WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&role, st)
	if err != nil {
		return role, err
	}

	return role, err
}

// GetRoleBySlug account from repo by slug.
func (ar *AuthRepo) GetRoleBySlug(slug string) (role model.Role, err error) {
	st := `SELECT * FROM roles WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&role, st)

	return role, err
}

// GetRoleByName account from repo by slug.
func (ar *AuthRepo) GetRoleByName(name string) (role model.Role, err error) {
	st := `SELECT * FROM roles WHERE name = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, name)

	err = ar.DB.Get(&role, st)

	return role, err
}

// Update role
func (ar *AuthRepo) UpdateRole(role *model.Role, tx ...*sqlx.Tx) error {
	ref, err := ar.GetRole(role.ID)
	if err != nil {
		return fmt.Errorf("cannot retrieve reference account: %s", err.Error())
	}

	role.SetUpdateValues()

	var st strings.Builder
	pcu := false // previous column updated?

	st.WriteString("UPDATE roles SET ")

	if role.Name.String != ref.Name.String {
		st.WriteString(kbs.SQLStrUpd("name", "name"))
		pcu = true
	}

	if role.Description.String != ref.Description.String {
		st.WriteString(kbs.SQLStrUpd("description", "description"))
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

	_, err = t.NamedExec(st.String(), role)

	if local {
		ar.Log.Info("Transaction created by repo: committing")
		return t.Commit()
	}

	return nil
}

// DeleteRole account from repo by ID.
func (ar *AuthRepo) DeleteRole(id uuid.UUID, tx ...*sqlx.Tx) error {
	st := `DELETE FROM role WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
func (ar *AuthRepo) DeleteRoleBySlug(slug string, tx ...*sqlx.Tx) error {
	st := `DELETE FROM roles WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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

// ResourcePermission --------------------------------------------------------------------------------
// Create a ResourcePermission
func (ar *AuthRepo) CreateResourcePermission(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error {
	st := `INSERT INTO resource_permission (id, slug, resource_id, permission_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
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

// GetResourcePermissionsByResourceID account from repo by slug.
func (ar *AuthRepo) GetResourcePermissionsByResourceID(id uuid.UUID) (resourcePermission []model.ResourcePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE resource_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetResourcePermissionsByPermissionID account from repo by slug.
func (ar *AuthRepo) GetResourcePermissionsByPermissionID(id uuid.UUID) (resourcePermission []model.ResourcePermission, err error) {
	st := `SELECT * FROM resource_permissions WHERE permission_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
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
	st := `DELETE FROM resource_permissions WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
	st := `DELETE FROM resource_permissions WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
	st := `INSERT INTO role_permissions (id, slug, role_id, permission_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
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
	st := `SELECT * FROM role_permissions WHERE is_deleted IS NULL OR NOT is_deleted`

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetRolePermission account by ID.
func (ar *AuthRepo) GetRolePermission(id uuid.UUID) (resourcePermission model.RolePermission, err error) {
	st := `SELECT * FROM role_permissions WHERE id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
	st = fmt.Sprintf(st, id.String())

	err = ar.DB.Get(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetRolePermissionBySlug account from repo by slug.
func (ar *AuthRepo) GetRolePermissionBySlug(slug string) (resourcePermission model.RolePermission, err error) {
	st := `SELECT * FROM role_permissions WHERE slug = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`

	st = fmt.Sprintf(st, slug)

	err = ar.DB.Get(&resourcePermission, st)

	return resourcePermission, err
}

// GetRolePermissiontByRoleID account from repo by slug.
func (ar *AuthRepo) GetRolePermissionByRoleID(id uuid.UUID) (resourcePermission []model.RolePermission, err error) {
	st := `SELECT * FROM role_permissions WHERE role_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&resourcePermission, st)
	if err != nil {
		return resourcePermission, err
	}

	return resourcePermission, err
}

// GetRolePermissiontByRoleID account from repo by slug.
func (ar *AuthRepo) GetRolePermissionByPermissionID(id uuid.UUID) (resourcePermission []model.RolePermission, err error) {
	st := `SELECT * FROM role_permissions WHERE permission_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
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

	st.WriteString("UPDATE role_permissions SET ")

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
	st := `DELETE FROM role_permissions WHERE id = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
	st := `DELETE FROM role_permissions WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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
	st := `INSERT INTO account_roles (id, slug, account_id, role_id, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
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

// GetAccountRolesByAccountID.
func (ar *AuthRepo) GetAccountRolesByAccountID(id uuid.UUID) (accountRole []model.AccountRole, err error) {
	st := `SELECT * FROM account_roles WHERE account_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted);`
	st = fmt.Sprintf(st, id)

	err = ar.DB.Select(&accountRole, st)
	if err != nil {
		return accountRole, err
	}

	return accountRole, err
}

// GetAccountRolesByRoleID.
func (ar *AuthRepo) GetAccountRolesByRoleID(id uuid.UUID) (accountRole []model.AccountRole, err error) {
	st := `SELECT * FROM account_roles WHERE role_id = '%s' AND (is_deleted IS NULL OR NOT is_deleted) LIMIT 1;`
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
	st := `DELETE FROM account_roles WHERE slug = '%s' AND (is_deleted IS NULL or NOT is_deleted);`
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

func (ar *AuthRepo) DeleteAccountRoleBySlugs(accountSlug, roleSlug string, tx ...*sqlx.Tx) error {
	st := `DELETE from account_roles
					WHERE account_roles.id IN (
						SELECT account_roles.id FROM account_roles
           		INNER JOIN accounts ON accounts.id = account_roles.account_id
          	  INNER JOIN roles ON roles.id = account_roles.role_id
           	WHERE accounts.slug = '%s'
					 		AND roles.slug = '%s'
             	AND (accounts.is_deleted IS NULL OR NOT accounts.is_deleted)
             	AND (accounts.is_active IS NULL OR accounts.is_active)
             	AND (roles.is_deleted IS NULL OR NOT roles.is_deleted)
             	AND (roles.is_active IS NULL OR roles.is_active)
             	AND (account_roles.is_deleted IS NULL OR NOT account_roles.is_deleted)
             	AND (account_roles.is_active IS NULL OR account_roles.is_active)
					);`

	st = fmt.Sprintf(st, accountSlug, roleSlug)

	t, local, err := ar.getTx(tx)
	if err != nil {
		ar.Log.Error(err)
		return err
	}

	_, err = t.Exec(st)
	if err != nil {
		return err
	}

	if local {
		err := t.Commit()
		return err
	}

	return err
}

// Custom

// GetAccountRoles
func (ar *AuthRepo) GetAccountRoles(accountSlug string) (roles []model.Role, err error) {
	st := `SELECT roles.* FROM roles
					INNER JOIN account_roles ON roles.id = account_roles.role_id
					INNER JOIN accounts ON accounts.id = account_roles.account_id
					WHERE accounts.slug = '%s'
						AND (accounts.is_deleted IS NULL OR NOT accounts.is_deleted)
						AND (accounts.is_active IS NULL OR accounts.is_active)
						AND (roles.is_deleted IS NULL OR NOT roles.is_deleted)
						AND (roles.is_active IS NULL OR roles.is_active)
						AND (account_roles.is_deleted IS NULL OR NOT account_roles.is_deleted)
						AND (account_roles.is_active IS NULL OR account_roles.is_active)
					ORDER BY roles.name ASC;`

	st = fmt.Sprintf(st, accountSlug)

	err = ar.DB.Select(&roles, st)
	if err != nil {
		return roles, err
	}

	return roles, err
}

func (ar *AuthRepo) GetNotAccountRoles(accountSlug string) (roles []model.Role, err error) {
	st := `SELECT roles.* from ROLES
					WHERE roles.id NOT IN (
						SELECT roles.id FROM roles
							INNER JOIN account_roles ON roles.id = account_roles.role_id
							INNER JOIN accounts ON accounts.id = account_roles.account_id
						WHERE accounts.slug = '%s'
							AND (accounts.is_deleted IS NULL OR NOT accounts.is_deleted)
							AND (accounts.is_active IS NULL OR accounts.is_active)
							AND (roles.is_deleted IS NULL OR NOT roles.is_deleted)
							AND (roles.is_active IS NULL OR roles.is_active)
							AND (account_roles.is_deleted IS NULL OR NOT account_roles.is_deleted)
							AND (account_roles.is_active IS NULL OR account_roles.is_active)
					);`

	st = fmt.Sprintf(st, accountSlug)

	err = ar.DB.Select(&roles, st)
	if err != nil {
		return roles, err
	}

	return roles, err
}

func (ar *AuthRepo) RemoveAccountRole(accountSlug, roleSlug string) (err error) {
	panic("not implemented")
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
