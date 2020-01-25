package mig

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
	// Migrator is a migrator handler.
	Migrator struct {
		*kbs.Migrator
	}
)

// NewMigrator creates and returns a new migrator.
func NewMigrator(cfg *kbs.Config, log kbs.Logger, name string, db *sqlx.DB) (*Migrator, error) {
	log.Info("New migrator", "name", name)

	m := &Migrator{
		kbs.NewMigrator(cfg, log, name, db),
	}

	m.addSteps()

	return m, nil
}

// GetMigrator configured.
func (m *Migrator) addSteps() {
	// Migrations
	// Enable Postgis
	s := &step{}
	s.Config(s.EnablePostgis, s.DropPostgis)
	m.AddMigration(s)

	// CreateUsersTable
	s = &step{}
	s.Config(s.CreateUsersTable, s.DropUsersTable)
	m.AddMigration(s)

	// CreateAccountsTable
	s = &step{}
	s.Config(s.CreateAccountsTable, s.DropAccountsTable)
	m.AddMigration(s)

	// CreateResourcesTable
	s = &step{}
	s.Config(s.CreateResourcesTable, s.DropResourcesTable)
	m.AddMigration(s)

	// CreateRolesTable
	s = &step{}
	s.Config(s.CreateRolesTable, s.DropRolesTable)
	m.AddMigration(s)

	// CreatePermissionsTable
	s = &step{}
	s.Config(s.CreatePermissionsTable, s.DropPermissionsTable)
	m.AddMigration(s)

	// CreateAccountRolesTable
	s = &step{}
	s.Config(s.CreateAccountRolesTable, s.DropAccountRolesTable)
	m.AddMigration(s)
}
