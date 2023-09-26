package ent

import (
	"context"
	"encoding/json"
	"entgo.io/ent/dialect"
	"fmt"
	"social/ent"
	"social/internal/app/http/middleware"
	"social/pkg/log"

	_ "github.com/lib/pq"
)

type PgConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Sslmode  string
}

func InitEnt(cfg PgConfig) *ent.Client {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Sslmode,
	)

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatal(err)
	}

	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			meta, ok := ctx.Value(middleware.AuditID).(middleware.Event)
			if !ok {
				return next.Mutate(ctx, mutation)
			}

			val, err := next.Mutate(ctx, mutation)

			meta.Table = mutation.Type()
			meta.Action = middleware.Action(mutation.Op().String())

			newValues, _ := json.Marshal(val)
			meta.NewValues = string(newValues)
			log.Info(meta)

			return val, err
		})
	})

	return client
}
