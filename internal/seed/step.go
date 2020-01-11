package seed

import (
	"github.com/jmoiron/sqlx"
	kbs "gitlab.com/kabestan/backend/kabestan"
)

type (
	step struct {
		name string
		seed kbs.SeedFx
		tx   *sqlx.Tx
	}
)

func (s *step) Config(seed kbs.SeedFx) {
	s.seed = seed
}

func (s *step) GetSeed() (up kbs.MigFx) {
	return s.seed
}

func (s *step) SetTx(tx *sqlx.Tx) {
	s.tx = tx
}

func (s *step) GetTx() (tx *sqlx.Tx) {
	return s.tx
}
