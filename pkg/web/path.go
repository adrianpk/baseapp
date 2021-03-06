package web

import "html/template"

var pathFxs = template.FuncMap{
	// User
	"userPath":           UserPath,
	"userPathEdit":       UserPathEdit,
	"userPathSlug":       UserPathSlug,
	"userPathInitDelete": UserPathInitDelete,
	"userPathNew":        UserPathNew,
	// Account
	"accountPath":           AccountPath,
	"accountPathEdit":       AccountPathEdit,
	"accountPathSlug":       AccountPathSlug,
	"accountPathInitDelete": AccountPathInitDelete,
	"accountPathNew":        AccountPathNew,
	"accountPathRoles":      AccountPathRoles,
	"accountPathRole":       AccountPathRole,
	// Resource
	"resourcePath":            ResourcePath,
	"resourcePathEdit":        ResourcePathEdit,
	"resourcePathSlug":        ResourcePathSlug,
	"resourcePathInitDelete":  ResourcePathInitDelete,
	"resourcePathNew":         ResourcePathNew,
	"resourcePathPermissions": ResourcePathPermissions,
	"resourcePathPermission":  ResourcePathPermission,
	// Role
	"rolePath":            RolePath,
	"rolePathEdit":        RolePathEdit,
	"rolePathSlug":        RolePathSlug,
	"rolePathInitDelete":  RolePathInitDelete,
	"rolePathNew":         RolePathNew,
	"rolePathPermissions": RolePathPermissions,
	"rolePathPermission":  RolePathPermission,
	// Path
	"permissionPath":           PermissionPath,
	"permissionPathEdit":       PermissionPathEdit,
	"permissionPathSlug":       PermissionPathSlug,
	"permissionPathInitDelete": PermissionPathInitDelete,
	"permissionPathNew":        PermissionPathNew,
}
