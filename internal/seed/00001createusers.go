package seed

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	zeroUUID   = "00000000-0000-0000-0000-000000000000"
	systemUUID = "00000000-0000-0000-0000-000000000001"
)

const (
	systemTenant      = "system"
	systemAccountType = "system"
)

const (
	utcTZ = "UTC"
	cetTZ = "CET"
)

var (
	// TODO: Create a builder for users that reads values from somewere: file, csv, map, etc...
	users = []map[string]interface{}{
		newUserMap("00000000-0000-0000-0000-000000000001", "system-000000000001", "system", "system", "system@kabestan.localhost"),

		newUserMap("00000000-0000-0000-0000-000000000002", "superadmin-000000000002", "superadmin", "superadmin", "superadmin@kabestan.localhost"),

		newUserMap("00000000-0000-0000-0000-000000000003", "admin-000000000003", "admin", "admin", "admin@kabestan.localhost"),
	}

	accounts = []map[string]interface{}{
		newAccountMap(users[0], systemAccountType, "", "", ""),

		newAccountMap(users[1], systemAccountType, "", "", ""),
	}
)

// CreateUsers seeding
func (s *step) CreateUsersAndAccounts() error {
	tx := s.GetTx()

	st := `INSERT INTO users (id, slug, username, password_digest, email, last_ip, confirmation_token, is_confirmed, geolocation, starts_at, ends_at, is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :username, :password_digest, :email, :last_ip, :confirmation_token, :is_confirmed, :geolocation, :starts_at, :ends_at, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// NOTE: Continue processing following after error?
	for _, u := range users {
		_, err := tx.NamedExec(st, u)
		if err != nil {
			log.Fatal(err)
		}
	}

	st = `INSERT INTO accounts (id, slug, tenant_id,  owner_id, parent_id, account_type, username, email, given_name, middle_names, family_name, locale, base_tz, current_tz,is_active, is_deleted, created_by_id, updated_by_id, created_at, updated_at)
VALUES (:id, :slug, :tenant_id, :owner_id, :parent_id, :account_type, :username, :email, :given_name, :middle_names, :family_name, :locale, :base_tz, :current_tz, :is_active, :is_deleted, :created_by_id, :updated_by_id, :created_at, :updated_at)`

	// NOTE: Continue processing following after error?
	for _, a := range accounts {
		_, err := tx.NamedExec(st, a)
		if err != nil {
			log.Fatal(err)
		}
	}

	return nil
}

func newUserMap(id, slug, username, password, email string) map[string]interface{} {
	return map[string]interface{}{
		"id":                 id,           //genUUID()
		"slug":               slug,         //genSlug(username),
		"tenant_id":          systemTenant, //genSlug(username),
		"username":           username,
		"password_digest":    genPassDigest(password),
		"email":              email,
		"last_ip":            "198.24.10.0/24",
		"confirmation_token": genUUIDStr(),
		"is_confirmed":       true,
		"geolocation":        "POINT(0 0)",
		"starts_at":          time.Now(),
		"ends_at":            time.Time{},
		"is_active":          true,
		"is_deleted":         false,
		"created_by_id":      systemUUID,
		"updated_by_id":      zeroUUID,
		"created_at":         time.Now(),
		"updated_at":         time.Time{},
	}
}

func genPassDigest(password string) string {
	pd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pd)
}

func genUUIDStr() string {
	return genUUID().String()
}

func genUUID() uuid.UUID {
	return uuid.NewV4()
}

func genSlug(prefix string) (slug string) {
	if strings.TrimSpace(prefix) == "" {
		prefix = "slug"
	}

	prefix = strings.Replace(prefix, "-", "", -1)
	prefix = strings.Replace(prefix, " ", "", -1)

	if !utf8.ValidString(prefix) {
		v := make([]rune, 0, len(prefix))
		for i, r := range prefix {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(prefix[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		prefix = string(v)
	}

	prefix = strings.ToLower(prefix)

	s := strings.Split(uuid.NewV4().String(), "-")
	l := s[len(s)-1]

	return strings.ToLower(fmt.Sprintf("%s-%s", prefix, l))
}

func newAccountMap(userMap map[string]interface{}, accountType, givenName, middleNames, familyName string) map[string]interface{} {
	return map[string]interface{}{
		"id":            userMap["id"],   //genUUID()
		"slug":          userMap["slug"], //genSlug(username),
		"tenant_id":     systemTenant,    //genSlug(username),
		"owner_id":      userMap["id"],
		"parent_id":     zeroUUID,
		"account_type":  accountType,
		"username":      userMap["username"],
		"email":         userMap["email"],
		"given_name":    givenName,
		"middle_names":  middleNames,
		"family_name":   familyName,
		"locale":        "en-US",
		"base_tz":       cetTZ,
		"current_tz":    cetTZ,
		"is_active":     true,
		"is_deleted":    false,
		"created_by_id": systemUUID,
		"updated_by_id": zeroUUID,
		"created_at":    time.Now(),
		"updated_at":    time.Time{},
	}
}
