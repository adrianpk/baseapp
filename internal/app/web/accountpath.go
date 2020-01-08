package web

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
)

// AccountRoot - Account resource root path.
var AccountRoot = "accounts"

// AccountPath
func AccountPath() string {
	return kbs.ResPath(AccountRoot)
}

// AccountPathEdit
func AccountPathEdit(res kbs.Identifiable) string {
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return kbs.ResPathEdit(AccountRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", AccountRoot, res.U)
}

// AccountPathNew
func AccountPathNew() string {
	return kbs.ResPathNew(AccountRoot)
}

// AccountPathInitDelete
func AccountPathInitDelete(res kbs.Identifiable) string {
	return kbs.ResPathInitDelete(AccountRoot, res)
}

// AccountPathSlug
func AccountPathSlug(res kbs.Identifiable) string {
	return kbs.ResPathSlug(AccountRoot, res)
}
