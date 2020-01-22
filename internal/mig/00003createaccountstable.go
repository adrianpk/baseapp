package mig

import "log"

// CreateAccountsTable migration
func (s *step) CreateAccountsTable() error {
	tx := s.GetTx()

	st := `CREATE TABLE accounts
		(
		id UUID PRIMARY KEY,
		slug VARCHAR(36) UNIQUE,
		tenant_id VARCHAR(128),
		owner_id UUID,
		parent_id UUID,
		account_type VARCHAR(36),
		username VARCHAR(32) UNIQUE,
		email VARCHAR(255),
		given_name VARCHAR(32),
		middle_names VARCHAR(32) NULL,
		family_name VARCHAR(64)
		);`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	st = `
		ALTER TABLE accounts
		ADD COLUMN locale VARCHAR(32),
		ADD COLUMN base_tz VARCHAR(64),
		ADD COLUMN current_tz VARCHAR(64),
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

// DropAccountsTable rollback
func (s *step) DropAccountsTable() error {
	tx := s.GetTx()

	st := `DROP TABLE accounts;`

	_, err := tx.Exec(st)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
