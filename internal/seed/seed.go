package seed

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // package init.
	kbs "gitlab.com/kabestan/backend/kabestan"
)

const (
	devDb  = "kabestan"
	testDb = "kabestan_test"
	prodDb = "kabestan_prod"
)

type (
	// Seeder is a seeder handler.
	Seeder struct {
		*kbs.Seeder
	}
)

// NewSeeder creates and returns a new seeder.
func NewSeeder(cfg *kbs.Config, log kbs.Logger, name string, db *sqlx.DB) (*Seeder, error) {
	log.Info("New seeder", "name", name)

	m := &Seeder{
		kbs.NewSeeder(cfg, log, name, db),
	}

	m.addSteps()

	return m, nil
}

// GetSeeder configured.
func (s *Seeder) addSteps() {
	// Seeds
	// Create user & accounts
	st := &step{}
	st.Config(st.CreateUsersAndAccounts)
	s.AddSeed(st)

	// Create resources
	st = &step{}
	st.Config(st.CreateResources)
	s.AddSeed(st)

	// Create roles
	st = &step{}
	st.Config(st.CreateRoles)
	s.AddSeed(st)

	// Create permissions
	st = &step{}
	st.Config(st.CreatePermissions)
	s.AddSeed(st)

	// Create account-roles
	st = &step{}
	st.Config(st.CreateAccountRoles)
	s.AddSeed(st)
}
