package main

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"social/config"
	"social/database"

	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	config.LoadConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.sslmode"),
	)
	migrator := database.Migrator(database.WithDSN(dsn))
	var err error
	migrator.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	migrator.Up()
}
