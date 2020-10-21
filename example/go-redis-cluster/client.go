package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redisext"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
)

var tracer = global.Tracer("go-redis-cluster-tracer")

func main() {
	ctx := context.Background()

	upclient := setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{"redis-node-0:6379", "redis-node-1:6379", "redis-node-2:6379"},

		NewClient: func(opt *redis.Options) *redis.Client {
			node := redis.NewClient(opt)
			node.AddHook(redisext.OpenTelemetryHook{})
			return node
		},
	})
	defer rdb.Close()

	rdb.AddHook(redisext.OpenTelemetryHook{})

	ctx, span := tracer.Start(ctx, "redis-main-span")
	defer span.End()

	if err := redisCommands(ctx, rdb); err != nil {
		upclient.ReportError(ctx, err)
		log.Println(err.Error())
		return
	}

	fmt.Println("trace", upclient.TraceURL(span))
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",
	})

	return upclient
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
