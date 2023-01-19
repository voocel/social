package database

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"testing"
)

func TestMigrator(t *testing.T) {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		"postgres",
		"postgres",
		"123456",
		"127.0.0.1",
		5433,
		"social",
		"disable",
	)
	migrator := Migrator(WithDSN(dsn))
	var err error
	migrator.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	migrator.Up()
}
