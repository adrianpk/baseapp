package web

import (
	"fmt"

	kbs "gitlab.com/kabestan/backend/kabestan"
)

// ResourceRoot - Resource resource root path.
var ResourceRoot = "resources"

// ResourcePath
func ResourcePath() string {
	return kbs.ResPath(ResourceRoot)
}

// ResourcePathEdit
func ResourcePathEdit(res kbs.Identifiable) string {
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return kbs.ResPathEdit(ResourceRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", ResourceRoot, res.U)
}

// ResourcePathNew
func ResourcePathNew() string {
	return kbs.ResPathNew(ResourceRoot)
}

// ResourcePathInitDelete
func ResourcePathInitDelete(res kbs.Identifiable) string {
	return kbs.ResPathInitDelete(ResourceRoot, res)
}

// ResourcePathSlug
func ResourcePathSlug(res kbs.Identifiable) string {
	return kbs.ResPathSlug(ResourceRoot, res)
}

// Custom

// ResourcePathPermissions
func ResourcePathPermissions(res kbs.Identifiable) string {
	return fmt.Sprintf("%s/permissions", kbs.ResPathSlug(ResourceRoot, res))
}

func ResourcePathPermission(res kbs.Identifiable, subRes kbs.Identifiable) string {
	return fmt.Sprintf("%s/permissions/%s", kbs.ResPathSlug(ResourceRoot, res), subRes.GetSlug())
}
