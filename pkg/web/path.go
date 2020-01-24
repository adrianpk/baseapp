package web

import "html/template"

var pathFxs = template.FuncMap{
	// User
	"userPath":           UserPath,
	"userPathEdit":       UserPathEdit,
	"userPathSlug":       UserPathSlug,
	"userPathInitDelete": UserPathInitDelete,
	"userPathNew":        UserPathNew,
	//
	"resourcePath":           ResourcePath,
	"resourcePathEdit":       ResourcePathEdit,
	"resourcePathSlug":       ResourcePathSlug,
	"resourcePathInitDelete": ResourcePathInitDelete,
	"resourcePathNew":        ResourcePathNew,
	// Role
	"rolePath":           RolePath,
	"rolePathEdit":       RolePathEdit,
	"rolePathSlug":       RolePathSlug,
	"rolePathInitDelete": RolePathInitDelete,
	"rolePathNew":        RolePathNew,
	// Path
	"permissionPath":           PermissionPath,
	"permissionPathEdit":       PermissionPathEdit,
	"permissionPathSlug":       PermissionPathSlug,
	"permissionPathInitDelete": PermissionPathInitDelete,
	"permissionPathNew":        PermissionPathNew,
}
