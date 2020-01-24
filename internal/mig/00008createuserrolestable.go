package mig

import "log"

// CreateUserRolesTable migration
func (s *step) CreateUserRolesTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE user_roles
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(36) UNIQUE,
		tenant_id VARCHAR(128),
		name VARCHAR(32) UNIQUE,
		user_id UUID NOT NULL,
		role_id UUID NOT NULL
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE user_roles
		ADD COLUMN is_active BOOLEAN,
		ADD COLUMN is_deleted BOOLEAN,
		ADD COLUMN created_by_id UUID,
		ADD COLUMN updated_by_id UUID,
		ADD COLUMN created_at TIMESTAMP,
		ADD COLUMN updated_at TIMESTAMP;`

	_, err = tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropUserRolesTable rollback
func (s *step) DropUserRolesTable() error {
	tx := s.GetTx()

	st := `DROP TABLE user_roles;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
