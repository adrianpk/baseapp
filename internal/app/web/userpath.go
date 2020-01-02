package web

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
)

// UserRoot - User resource root path.
var UserRoot = "users"

// UserPath
func UserPath() string {
	return kbs.ResPath(UserRoot)
}

// UserPathEdit
func UserPathEdit(res kbs.Identifiable) string {
	// TODO: Analize if in a multi-tenant setup this could be
	// a problem.
	return kbs.ResPathEdit(UserRoot, res)
	//return fmt.Sprintf("/%s/%s/edit", UserRoot, res.U)
}

// UserPathNew
func UserPathNew() string {
	return kbs.ResPathNew(UserRoot)
}

// UserPathInitDelete
func UserPathInitDelete(res kbs.Identifiable) string {
	return kbs.ResPathInitDelete(UserRoot, res)
}

// UserPathSlug
func UserPathSlug(res kbs.Identifiable) string {
	return kbs.ResPathSlug(UserRoot, res)
}

// UserPathSignUp
func UserPathSignUp() string {
	return kbs.ResPath(UserRoot) + "/signup"
}

// UserPathSignIn
func UserPathSignIn() string {
	return kbs.ResPath(UserRoot) + "/signin"
}
