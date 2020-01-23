package mig

import "log"

// CreateRolesTable migration
func (s *step) CreateRolesTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE roles
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(36) UNIQUE,
		tenant_id VARCHAR(128),
		name VARCHAR(32) UNIQUE,
		description TEXT NULL,
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE roles
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

// DropRolesTable rollback
func (s *step) DropRolesTable() error {
	tx := s.GetTx()

	st := `DROP TABLE roles;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
