package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{":9123", ":9124", ":9125"},

		NewClient: func(opt *redis.Options) *redis.Client {
			node := redis.NewClient(opt)
			node.AddHook(&redisotel.TracingHook{})
			return node
		},
	})
	defer rdb.Close()

	rdb.AddHook(&redisotel.TracingHook{})

	ctx, span := tracer.Start(ctx, "redis-main-span")
	defer span.End()

	if err := redisCommands(ctx, rdb); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("trace", uptrace.TraceURL(span))
}

func redisCommands(ctx context.Context, rdb *redis.ClusterClient) error {
	if err := rdb.Set(ctx, "foo", "bar", 0).Err(); err != nil {
		return err
	}

	if err := rdb.Get(ctx, "foo").Err(); err != nil {
		return err
	}

	if _, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "foo", "bar2", 0)
		pipe.Get(ctx, "foo")
		return nil
	}); err != nil {
		return err
	}

	return nil
}
