package web

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
)

// PermissionRoot - Permission permission root path.
var PermissionRoot = "permissions"

// PermissionPath
func PermissionPath() string {
	return kbs.ResPath(PermissionRoot)
}

// PermissionPathEdit
func PermissionPathEdit(res kbs.Identifiable) string {
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return kbs.ResPathEdit(PermissionRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", PermissionRoot, res.U)
}

// PermissionPathNew
func PermissionPathNew() string {
	return kbs.ResPathNew(PermissionRoot)
}

// PermissionPathInitDelete
func PermissionPathInitDelete(res kbs.Identifiable) string {
	return kbs.ResPathInitDelete(PermissionRoot, res)
}

// PermissionPathSlug
func PermissionPathSlug(res kbs.Identifiable) string {
	return kbs.ResPathSlug(PermissionRoot, res)
}
