package mig

// EnablePostgis migration
func (s *step) EnablePostgis() error {
	tx := s.GetTx()

	st := `CREATE EXTENSION IF NOT EXISTS postgis;`
	//CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}

// DropPostgis rollback
func (s *step) DropPostgis() error {
	tx := s.GetTx()

	st := `DROP EXTENSION IF EXISTS postgis;`

	_, err := tx.Exec(st)
	if err != nil {
		return err
	}

	return nil
}
