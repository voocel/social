package database

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migrate struct {
	DB  *sql.DB
	dsn string
}

type Options func(opts *Migrate)

func Migrator(opts ...Options) *Migrate {
	m := &Migrate{}
	goose.SetBaseFS(embedMigrations)

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Migrate) Up() {
	if err := goose.Up(m.DB, "migrations"); err != nil {
		panic(err)
	}
}

func (m *Migrate) Down() {
	if err := goose.Down(m.DB, "migrations"); err != nil {
		panic(err)
	}
}

func WithDSN(dsn string) func(opts *Migrate) {
	return func(opts *Migrate) {
		opts.dsn = dsn
	}
}
