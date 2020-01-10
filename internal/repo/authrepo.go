package repo

import (
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/kabestan/repo/baseapp/internal/model"
)

type (
	AuthRepo interface {
		// Resource
		CreateResource(u *model.Resource, tx ...*sqlx.Tx) error
		GetAllResources() (resources []model.Resource, err error)
		GetResource(id uuid.UUID) (resource model.Resource, err error)
		GetResourceBySlug(slug string) (resource model.Resource, err error)
		GetResourceByName(name string) (model.Resource, error)
		GetResourceByTag(tag string) (resource model.Resource, err error)
		GetResourceByPath(path string) (resource model.Resource, err error)
		UpdateResource(resource *model.Resource, tx ...*sqlx.Tx) error
		DeleteResource(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteResourceBySlug(slug string, tx ...*sqlx.Tx) error
		// Permission
		CreatePermission(u *model.Permission, tx ...*sqlx.Tx) error
		GetAllPermissions() (permissions []model.Permission, err error)
		GetPermission(id uuid.UUID) (permission model.Permission, err error)
		GetPermissionBySlug(slug string) (permission model.Permission, err error)
		GetPermissionByName(name string) (model.Permission, error)
		UpdatePermission(permission *model.Permission, tx ...*sqlx.Tx) error
		DeletePermission(id uuid.UUID, tx ...*sqlx.Tx) error
		DeletePermissionBySlug(slug string, tx ...*sqlx.Tx) error
		// Role
		CreateRole(u *model.Role, tx ...*sqlx.Tx) error
		GetAllRoles() (roles []model.Role, err error)
		GetRole(id uuid.UUID) (role model.Role, err error)
		GetRoleBySlug(slug string) (role model.Role, err error)
		GetRoleByName(name string) (model.Role, error)
		UpdateRole(role *model.Role, tx ...*sqlx.Tx) error
		DeleteRole(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteRoleBySlug(slug string, tx ...*sqlx.Tx) error
		// ResourcePermission
		CreateResourcePermission(u *model.ResourcePermission, tx ...*sqlx.Tx) error
		GetAllResourcePermissions() (resourcePermissions []model.ResourcePermission, err error)
		GetResourcePermission(id uuid.UUID) (resourcePermission model.ResourcePermission, err error)
		GetResourcePermissionBySlug(slug string) (resourcePermission model.ResourcePermission, err error)
		GetResourcePermissionsByResourceID(uuid.UUID) ([]model.ResourcePermission, error)
		GetResourcePermissionsByPermissionID(uuid.UUID) ([]model.ResourcePermission, error)
		UpdateResourcePermission(resourcePermission *model.ResourcePermission, tx ...*sqlx.Tx) error
		DeleteResourcePermission(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteResourcePermissionBySlug(slug string, tx ...*sqlx.Tx) error
		// RolePermission
		CreateRolePermission(u *model.RolePermission, tx ...*sqlx.Tx) error
		GetAllRolePermissions() (rolePermissions []model.RolePermission, err error)
		GetRolePermission(id uuid.UUID) (rolePermission model.RolePermission, err error)
		GetRolePermissionBySlug(slug string) (rolePermission model.RolePermission, err error)
		GetRolePermissionByRoleID(uuid.UUID) ([]model.RolePermission, error)
		GetRolePermissionByPermissionID(uuid.UUID) ([]model.RolePermission, error)
		UpdateRolePermission(rolePermission *model.RolePermission, tx ...*sqlx.Tx) error
		DeleteRolePermission(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteRolePermissionBySlug(slug string, tx ...*sqlx.Tx) error
		// AccountRole
		CreateAccountRole(u *model.AccountRole, tx ...*sqlx.Tx) error
		GetAllAccountRoles() (auth []model.AccountRole, err error)
		GetAccountRole(id uuid.UUID) (auth model.AccountRole, err error)
		GetAccountRoleBySlug(slug string) (auth model.AccountRole, err error)
		GetAccountRoleByAccountID(uuid.UUID) ([]model.AccountRole, error)
		GetAccountRoleByRoleID(uuid.UUID) ([]model.AccountRole, error)
		UpdateAccountRole(auth *model.AccountRole, tx ...*sqlx.Tx) error
		DeleteAccountRole(id uuid.UUID, tx ...*sqlx.Tx) error
		DeleteAccountRoleBySlug(slug string, tx ...*sqlx.Tx) error
	}
)
