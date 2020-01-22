package kabestan

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type (
	SeederIF interface {
		Seed() error
	}
)

type (
	// Fx type alias
	SeedFx = func() error

	// Seeder struct.
	Seeder struct {
		*Worker
		DB     *sqlx.DB
		schema string
		dbName string
		seeds  []*Seed
	}

	// Exec interface.
	SeedExec interface {
		Config(seed SeedFx)
		GetSeed() (up SeedFx)
		SetTx(tx *sqlx.Tx)
		GetTx() (tx *sqlx.Tx)
	}

	// Seed struct.
	Seed struct {
		Executor SeedExec
	}

	seedRecord struct {
		ID        uuid.UUID `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		Fx        string    `db:"up_fx" json:"upFx"`
		IsApplied bool      `db:"is_applied" json:"isApplied"`
		CreatedAt time.Time `db:"created_at" json:"createdAt"`
	}
)

const (
	pgSeederTable = "seeds"

	pgCreateSeederSt = `CREATE TABLE %s.%s (
		id UUID PRIMARY KEY,
		name VARCHAR(64),
		fx VARCHAR(64),
 		is_applied BOOLEAN,
		created_at TIMESTAMP
	);`

	pgDropSeederSt = `DROP TABLE %s.%s;`

	pgSelSeederSt = `SELECT is_applied FROM %s.%s WHERE name = '%s' and is_applied = true;`

	pgRecSeederSt = `INSERT INTO %s.%s (id, name, fx, is_applied, created_at)
		VALUES (:id, :name, :fx, :is_applied, :created_at);`

	pgDelSeederSt = `DELETE FROM %s.%s WHERE name = '%s' and is_applied = true;`
)

// NewSeeder.
func NewSeeder(cfg *Config, log Logger, name string, db *sqlx.DB) *Seeder {
	m := &Seeder{
		Worker: NewWorker(cfg, log, name),
		DB:     db,
		schema: cfg.ValOrDef("pg.schema", ""),
		dbName: cfg.ValOrDef("pg.database", ""),
	}

	return m
}

// pgConnect to postgre database
// mainly user to create and drop app database.
func (m *Seeder) pgConnect() error {
	db, err := sqlx.Open("postgres", m.pgDbURL())
	if err != nil {
		log.Printf("Connection error: %s\n", err.Error())
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Connection error: %s", err.Error())
		return err
	}

	m.DB = db
	return nil
}

// GetTx returns a new transaction from seeder connection.
func (s *Seeder) GetTx() *sqlx.Tx {
	return s.DB.MustBegin()
}

// PreSetup creates database
// and seeder table if needed.
func (s *Seeder) PreSetup() {
	if !s.dbExists() {
		s.CreateDb()
	}

	if !s.seedTableExists() {
		s.createSeederTable()
	}
}

// dbExists returns true if seeder
// referenced database has been already created.
// Only for postgress at the moment.
func (s *Seeder) dbExists() bool {
	st := fmt.Sprintf(`SELECT EXISTS(
		SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s'));`, s.dbName)

	r, err := s.DB.Query(st)
	if err != nil {
		log.Printf("Error checking database: %s\n", err.Error())
		return false
	}

	for r.Next() {
		var exists sql.NullBool
		err = r.Scan(&exists)
		if err != nil {
			log.Printf("Cannot read query result: %s\n", err.Error())
			return false
		}
		return exists.Bool
	}
	return false
}

// seedExists returns true if seeder table exists.
func (s *Seeder) seedTableExists() bool {
	st := fmt.Sprintf(`SELECT EXISTS (
		SELECT 1
   	FROM   pg_catalog.pg_class c
   	JOIN   pg_catalog.pg_namespace n ON n.oid = c.relnamespace
   	WHERE  n.nspname = '%s'
   	AND    c.relname = '%s'
   	AND    c.relkind = 'r'
	);`, s.schema, s.dbName)

	r, err := s.DB.Query(st)
	if err != nil {
		log.Printf("Error checking database: %s\n", err.Error())
		return false
	}

	for r.Next() {
		var exists sql.NullBool
		err = r.Scan(&exists)
		if err != nil {
			log.Printf("Cannot read query result: %s\n", err.Error())
			return false
		}

		return exists.Bool
	}
	return false
}

// CreateDb for seeder.
func (s *Seeder) CreateDb() (string, error) {
	//s.CloseAppConns()
	st := fmt.Sprintf(pgCreateDbSt, s.dbName)

	_, err := s.DB.Exec(st)
	if err != nil {
		return s.dbName, err
	}

	return s.dbName, nil
}

func (s *Seeder) createSeederTable() (string, error) {
	tx := s.GetTx()

	st := fmt.Sprintf(pgCreateSeederSt, s.schema, pgSeederTable)

	_, err := tx.Exec(st)
	if err != nil {
		return pgSeederTable, err
	}

	return pgSeederTable, tx.Commit()
}

func (s *Seeder) AddSeed(e SeedExec) {
	s.seeds = append(s.seeds, &Seed{Executor: e})
}

func (s *Seeder) Seed() error {
	s.PreSetup()

	for _, sd := range s.seeds {
		fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaa")
		fmt.Printf("%+v/n/n", *sd)
		exec := sd.Executor
		fn := getFxName(exec.GetSeed())
		name := seedName(fn)

		// Continue if already applied
		if !s.canApplySeed(name) {
			log.Printf("Seed '%s' already applied.", name)
			continue
		}

		// Get a new Tx from seeder
		tx := s.GetTx()
		// Pass Tx to the executor
		exec.SetTx(tx)

		// Execute seed
		values := reflect.ValueOf(exec).MethodByName(fn).Call([]reflect.Value{})

		// Read error
		err, ok := values[0].Interface().(error)
		if !ok && err != nil {
			log.Printf("Seed step not executed: %s\n", fn) // TODO: Remove log
			log.Printf("Err  %+v' of type %T\n", err, err) // TODO: Remove log.
			msg := fmt.Sprintf("cannot run seeding '%s': %s", fn, err.Error())
			tx.Rollback()
			return errors.New(msg)
		}
		fmt.Println("aaaaaaaaaaaaaaaa")
		// Register seed
		err = s.recSeed(exec)
		fmt.Println("bbbbbbbbbbbbbbbb")

		err = tx.Commit()
		if err != nil {
			msg := fmt.Sprintf("Commit error: %s\n", err.Error())
			log.Printf("Commit error: %s", msg)
			tx.Rollback()
			return errors.New(msg)
		}

		log.Printf("Seed step executed: %s\n", fn)
	}

	return nil
}

func (s *Seeder) canApplySeed(name string) bool {
	st := fmt.Sprintf(pgSelSeederSt, s.schema, pgSeederTable, name)
	r, err := s.DB.Query(st)

	if err != nil {
		log.Printf("Cannot determine seeder status: %s\n", err.Error())
		return false
	}

	for r.Next() {
		var applied sql.NullBool
		err = r.Scan(&applied)
		if err != nil {
			log.Printf("Cannot determine seeder status: %s\n", err.Error())
			return false
		}

		return !applied.Bool
	}

	return true
}

func (s *Seeder) recSeed(e SeedExec) error {
	log.Println("1:----------------------------")
	st := fmt.Sprintf(pgRecSeederSt, s.schema, pgSeederTable)
	fx := getFxName(e.GetSeed())
	name := seedName(fx)
	log.Printf("%+s", fx)

	_, err := e.GetTx().NamedExec(st, seedRecord{
		ID:        uuid.NewV4(),
		Name:      name,
		Fx:        fx,
		IsApplied: true,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("2:----------------------------")
		msg := fmt.Sprintf("Cannot update seeder table: %s\n", err.Error())
		return errors.New(msg)
	}

	log.Println("3:----------------------------")
	return nil
}

func seedName(fxName string) string {
	return toSnakeCase(fxName)
}

func (m *Seeder) dbURL() string {
	host := m.Cfg.ValOrDef("pg.host", "localhost")
	port := m.Cfg.ValOrDef("pg.port", "5432")
	m.schema = m.Cfg.ValOrDef("pg.schema", "public")
	m.dbName = m.Cfg.ValOrDef("pg.database", "kabestan_test_d1x89s0l")
	user := m.Cfg.ValOrDef("pg.user", "kabestan")
	pass := m.Cfg.ValOrDef("pg.password", "kabestan")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbName=%s sslmode=disable search_path=%s", host, port, user, pass, m.dbName, m.schema)
}

func (m *Seeder) pgDbURL() string {
	host := m.Cfg.ValOrDef("pg.host", "localhost")
	port := m.Cfg.ValOrDef("pg.port", "5432")
	schema := "public"
	db := "postgres"
	user := m.Cfg.ValOrDef("pg.user", "kabestan")
	pass := m.Cfg.ValOrDef("pg.password", "kabestan")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbName=%s sslmode=disable search_path=%s", host, port, user, pass, db, schema)
}
