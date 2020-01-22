package model

import (
	"database/sql"

	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	kbs "gitlab.com/kabestan/backend/kabestan"
	"gitlab.com/kabestan/backend/kabestan/db"
	"golang.org/x/crypto/bcrypt"
)

type (
	// User model
	User struct {
		kbs.Identification
		Username          sql.NullString `db:"username" json:"username" schema:"username"`
		Password          string         `db:"-" json:"password" schema:"password"`
		PasswordDigest    sql.NullString `db:"password_digest" json:"-" schema:"-"`
		Email             sql.NullString `db:"email" json:"email" schema:"email"`
		EmailConfirmation sql.NullString `db:"-" json:"emailConfirmation" schema:"email-confirmation"`
		LastIP            sql.NullString `db:"last_ip" json:"-" schema:"-"`
		ConfirmationToken sql.NullString `db:"confirmation_token" json:"-" schema:"-"`
		IsConfirmed       sql.NullBool   `db:"is_confirmed" json:"-" schema:"-"`
		Geolocation       kbs.Point      `db:"geolocation" json:"-" schema:"-"`
		StartsAt          pq.NullTime    `db:"starts_at" json:"-" schema:"-"`
		EndsAt            pq.NullTime    `db:"ends_at" json:"-" schema:"-"`
		IsActive          sql.NullBool   `db:"is_active" json:"-" schema:"-"`
		IsDeleted         sql.NullBool   `db:"is_deleted" json:"-" schema:"-"`
		kbs.Audit
	}
)

type (
	UserForm struct {
		Slug              string `json:"slug" schema:"slug"`
		Username          string `json:"username" schema:"username"`
		Password          string `json:"password" schema:"password"`
		Email             string `json:"email" schema:"email"`
		EmailConfirmation string `json:"emailConfirmation" schema:"email-confirmation"`
		IsNew             bool   `json:"-" schema:"-"`
	}
)

func ToUserFormList(users []User) (fs []UserForm) {
	for _, m := range users {
		fs = append(fs, m.ToForm())
	}
	return fs
}

// UpdatePasswordDigest if password changed.
func (user *User) UpdatePasswordDigest() (digest string, err error) {
	if user.Password == "" {
		return user.PasswordDigest.String, nil
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.PasswordDigest.String, err
	}
	user.PasswordDigest = db.ToNullString(string(hpass))
	return user.PasswordDigest.String, nil
}

// SetCreateValues sets de ID and slug.
func (user *User) SetCreateValues() error {
	// Set create values only only if they were not previously
	if user.Identification.ID == uuid.Nil ||
		user.Identification.Slug.String == "" {
		pfx := user.Username.String
		user.Identification.SetCreateValues(pfx)
		user.Audit.SetCreateValues()
		user.UpdatePasswordDigest()
	}
	return nil
}

// SetUpdateValues
func (user *User) SetUpdateValues() error {
	user.Audit.SetUpdateValues()
	user.UpdatePasswordDigest()
	return nil
}

// GenConfirmationToken
func (user *User) GenConfirmationToken() {
	user.ConfirmationToken = db.ToNullString(uuid.NewV4().String())
	user.IsConfirmed = db.ToNullBool(false)
}

// GenAutoConfirmationToken
func (user *User) GenAutoConfirmationToken() {
	user.ConfirmationToken = db.ToNullString(uuid.NewV4().String())
	user.IsConfirmed = db.ToNullBool(true)
}

// Match condition for
func (user *User) Match(tc *User) bool {
	r := user.Identification.Match(tc.Identification) &&
		user.Username == tc.Username &&
		user.PasswordDigest == tc.PasswordDigest &&
		user.Email == tc.Email &&
		user.IsConfirmed == tc.IsConfirmed &&
		user.Geolocation == tc.Geolocation &&
		user.StartsAt == tc.StartsAt &&
		user.EndsAt == tc.EndsAt
	return r
}

// ToForm lets convert a model to its associated form type.
// This convertion step could be avoided since gorilla schema allows
// to register custom decoders and in this case we need one because
// struct properties are not using Go standard types but their sql
// null conterpart types. As long as is relatively simple, ergonomic
// and could be easily implemented with generators I prefer to avoid
// the use of reflection.
func (user *User) ToForm() UserForm {
	return UserForm{
		Slug:              user.Slug.String,
		Username:          user.Username.String,
		Email:             user.Email.String,
		EmailConfirmation: user.Email.String,
		IsNew:             user.IsNew(),
	}
}

// ToModel lets covert a form type to its associated model.
func (userForm *UserForm) ToModel() User {
	return User{
		Identification: kbs.Identification{
			Slug: db.ToNullString(userForm.Slug),
		},
		Username:          db.ToNullString(userForm.Username),
		Password:          userForm.Password,
		Email:             db.ToNullString(userForm.Email),
		EmailConfirmation: db.ToNullString(userForm.EmailConfirmation),
	}
}

func (userForm *UserForm) GetSlug() string {
	return userForm.Slug
}
