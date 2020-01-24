package mig

import "log"

// CreatePermissionsTable migration
func (s *step) CreatePermissionsTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE permissions
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(36) UNIQUE,
		tenant_id VARCHAR(128),
		name VARCHAR(32) UNIQUE,
		description TEXT NULL,
		tag VARCHAR(16) UNIQUE,
		path VARCHAR(512) UNIQUE
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE permissions
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

// DropPermissionsTable rollback
func (s *step) DropPermissionsTable() error {
	tx := s.GetTx()

	st := `DROP TABLE permissions;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
