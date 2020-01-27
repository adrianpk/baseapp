package mig

import "log"

// CreateAccountRolesTable migration
func (s *step) CreateAccountRolesTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE account_roles
	(
		id UUID PRIMARY KEY,
		slug VARCHAR(64) UNIQUE,
		tenant_id VARCHAR(128),
		name VARCHAR(32) UNIQUE,
		account_id UUID NOT NULL,
		role_id UUID NOT NULL
	);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE account_roles
		ADD COLUMN is_active BOOLEAN,
		ADD COLUMN is_deleted BOOLEAN,
		ADD COLUMN created_by_id UUID,
		ADD COLUMN updated_by_id UUID,
		ADD COLUMN created_at TIMESTAMP,
		ADD COLUMN updated_at TIMESTAMP,
		ADD UNIQUE (tenant_id, account_id, role_id);`

	_, err = tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropAccountRolesTable rollback
func (s *step) DropAccountRolesTable() error {
	tx := s.GetTx()

	st := `DROP TABLE account_roles;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
