package main

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redisext"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
)

var tracer = global.Tracer("go-redis-tracer")

func main() {
	ctx := context.Background()

	upclient := setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-server:6379",
	})
	rdb.AddHook(&redisext.OpenTelemetryHook{})
	defer rdb.Close()

	err := traceRedisCommands(ctx, rdb)
	if err != nil {
		upclient.ReportError(ctx, err)
		log.Println(err.Error())
		return
	}
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	hostname, _ := os.Hostname()
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",

		Resource: map[string]interface{}{
			"host.name": hostname,
		},
	})

	return upclient
}

func traceRedisCommands(ctx context.Context, rdb *redis.Client) error {
	ctx, span := tracer.Start(ctx, "redis-commands")
	defer span.End()

	err := rdb.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		return err
	}

	err = rdb.Get(ctx, "foo").Err()
	if err != nil {
		return err
	}

	_, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "foo", "bar2", 0)
		pipe.Get(ctx, "foo")
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
