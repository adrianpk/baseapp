package mig

import "log"

// CreateResourcePermissionsTable migration
func (s *step) CreateResourcePermissionsTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE resource_permissions
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(64) UNIQUE,
		tenant_id VARCHAR(128),
		name VARCHAR(32) UNIQUE,
		resource_id UUID NOT NULL,
		permission_id UUID NOT NULL
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE resource_permissions
		ADD COLUMN is_active BOOLEAN,
		ADD COLUMN is_deleted BOOLEAN,
		ADD COLUMN created_by_id UUID,
		ADD COLUMN updated_by_id UUID,
		ADD COLUMN created_at TIMESTAMP,
		ADD COLUMN updated_at TIMESTAMP,
		ADD UNIQUE (tenant_id, resource_id, permission_id);`

	_, err = tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropResourcePermissionsTable rollback
func (s *step) DropResourcePermissionsTable() error {
	tx := s.GetTx()

	st := `DROP TABLE resource_permissions;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
