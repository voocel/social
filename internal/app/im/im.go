package im

import (
	"context"
	"encoding/json"
	"entgo.io/ent/dialect"
	"fmt"
	"github.com/spf13/viper"
	"social/ent"
	"social/internal/app/http/middleware"
	"social/internal/node"
	"social/internal/transport"
	"social/internal/transport/grpc"
	"social/pkg/log"

	_ "github.com/lib/pq"
)

func Run() *node.Node {
	addr := viper.GetString("transport.grpc.addr")
	name := viper.GetString("transport.grpc.service_name")
	n := node.NewNode(
		node.WithTransporter(
			grpc.NewTransporter(transport.WithAddr(addr), transport.WithName(name)),
		),
	)
	client := initEnt()
	core := newCore(n.GetProxy(), client)
	core.Init()
	n.Start()
	return n
}

func initEnt() *ent.Client {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.sslmode"),
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
