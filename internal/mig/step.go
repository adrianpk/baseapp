package mig

import (
	kbs "gitlab.com/kabestan/backend/kabestan"
	"github.com/jmoiron/sqlx"
)

type (
	step struct {
		name string
		up   kbs.MigFx
		down kbs.MigFx
		tx   *sqlx.Tx
	}
)

func (s *step) Config(up kbs.MigFx, down kbs.MigFx) {
	s.up = up
	s.down = down
}

func (s *step) GetName() (name string) {
	return s.name
}

func (s *step) GetUp() (up kbs.MigFx) {
	return s.up
}

func (s *step) GetDown() (down kbs.MigFx) {
	return s.down
}

func (s *step) SetTx(tx *sqlx.Tx) {
	s.tx = tx
}

func (s *step) GetTx() (tx *sqlx.Tx) {
	return s.tx
}
