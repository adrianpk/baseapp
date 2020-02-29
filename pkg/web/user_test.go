package web_test

import (
	"net/url"
	"testing"

	"gitlab.com/kabestan/repo/baseapp/pkg/web"
)

var (
	userDataValid = map[string]string{
		"username":          "username",
		"password":          "password",
		"email":             "username@mail.com",
		"emailConfirmation": "username@mail.com",
		"givenName":         "name",
		"middleNames":       "middles",
		"familyName":        "family",
	}

	userUpdateDataValid = map[string]string{
		"username":          "usernameUpd",
		"password":          "passwordUpd",
		"email":             "usernameUpd@mail.com",
		"emailConfirmation": "usernameUpd@mail.com",
		"givenName":         "nameUpd",
		"middleNames":       "middlesUpd",
		"familyName":        "familyUpd",
	}

	userSample1 = map[string]string{
		"username":          "username1",
		"password":          "password1",
		"email":             "username1@mail.com",
		"emailConfirmation": "username1@mail.com",
		"givenName":         "name1",
		"middleNames":       "middles1",
		"familyName":        "family1",
	}

	userSample2 = map[string]string{
		"username":          "username2",
		"password":          "password2",
		"email":             "username2@mail.com",
		"emailConfirmation": "username2@mail.com",
		"givenName":         "name2",
		"middleNames":       "middles2",
		"familyName":        "family2",
	}
)

// TestCreateUser tests user creation.
func TestCreateUser(t *testing.T) {
	t.Log("TestCreateUser init")
	// Make request
	userURL := buildURL(web.UserRoot)

	userForm := url.Values{
		"username":          []string{userDataValid["username"]},
		"Password":          []string{userDataValid["password"]},
		"email":             []string{userDataValid["email"]},
		"emailConfirmation": []string{userDataValid["emailConfirmation"]},
		"givenName":         []string{userDataValid["givenName"]},
		"middleNames":       []string{userDataValid["middleNames"]},
		"familyName":        []string{userDataValid["familyName"]},
	}

	clt := ts.Client
	body, err := clt.Post(userURL, userForm)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%+v", body)
	t.Log("TestCreateUser end")
}
