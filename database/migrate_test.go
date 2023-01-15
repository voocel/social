package database

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/copier"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
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

func TestName(t *testing.T) {
	type Tt struct {
		Name string
		Age  string
	}
	j := Tt{
		Name: "aaa",
		Age:  "10",
	}
	m := make(map[string]interface{})
	err := copier.Copy(&m, &j)
	fmt.Println(err)
	fmt.Printf("%+v\n", m)
}
