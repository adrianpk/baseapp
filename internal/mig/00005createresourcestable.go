package mig

import "log"

// CreateResourcesTable migration
func (s *step) CreateResourcesTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE resources
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(64) UNIQUE,
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
		ALTER TABLE resources
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

// DropResourcesTable rollback
func (s *step) DropResourcesTable() error {
	tx := s.GetTx()

	st := `DROP TABLE resources;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
