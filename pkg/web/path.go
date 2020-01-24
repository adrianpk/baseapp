package web

import "html/template"

var pathFxs = template.FuncMap{
	// User
	"userPath":           UserPath,
	"userPathEdit":       UserPathEdit,
	"userPathSlug":       UserPathSlug,
	"userPathInitDelete": UserPathInitDelete,
	"userPathNew":        UserPathNew,
	// Resource
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
}
