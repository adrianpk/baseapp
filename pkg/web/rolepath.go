package web

import (
	"fmt"

	kbs "gitlab.com/kabestan/backend/kabestan"
)

// RoleRoot - Role role root path.
var RoleRoot = "roles"

// RolePath
func RolePath() string {
	return kbs.ResPath(RoleRoot)
}

// RolePathEdit
func RolePathEdit(res kbs.Identifiable) string {
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return kbs.ResPathEdit(RoleRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", RoleRoot, res.U)
}

// RolePathNew
func RolePathNew() string {
	return kbs.ResPathNew(RoleRoot)
}

// RolePathInitDelete
func RolePathInitDelete(res kbs.Identifiable) string {
	return kbs.ResPathInitDelete(RoleRoot, res)
}

// RolePathSlug
func RolePathSlug(res kbs.Identifiable) string {
	return kbs.ResPathSlug(RoleRoot, res)
}

// Custom

// RolePathPermissions
func RolePathPermissions(res kbs.Identifiable) string {
	return fmt.Sprintf("%s/permissions", kbs.ResPathSlug(RoleRoot, res))
}

func RolePathPermission(res kbs.Identifiable, subRes kbs.Identifiable) string {
	return fmt.Sprintf("%s/permissions/%s", kbs.ResPathSlug(RoleRoot, res), subRes.GetSlug())
}
