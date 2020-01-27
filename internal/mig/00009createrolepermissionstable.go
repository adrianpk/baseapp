package mig

import "log"

// CreateRolePermissionsTable migration
func (s *step) CreateRolePermissionsTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE role_permissions
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(64) UNIQUE,
		tenant_id VARCHAR(128),
		name VARCHAR(32) UNIQUE,
		role_id UUID NOT NULL,
		permission_id UUID NOT NULL
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE role_permissions
		ADD COLUMN is_active BOOLEAN,
		ADD COLUMN is_deleted BOOLEAN,
		ADD COLUMN created_by_id UUID,
		ADD COLUMN updated_by_id UUID,
		ADD COLUMN created_at TIMESTAMP,
		ADD COLUMN updated_at TIMESTAMP,
		ADD UNIQUE (tenant_id, role_id, permission_id);`

	_, err = tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropRolePermissionsTable rollback
func (s *step) DropRolePermissionsTable() error {
	tx := s.GetTx()

	st := `DROP TABLE role_permissions;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
