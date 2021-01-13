package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("go-redis-tracer")

func main() {
	ctx := context.Background()

	upclient := newUptraceClient()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-server:6379",
	})
	defer rdb.Close()

	rdb.AddHook(&redisotel.TracingHook{})

	ctx, span := tracer.Start(ctx, "redis-main-span")
	defer span.End()

	if err := redisCommands(ctx, rdb); err != nil {
		upclient.ReportError(ctx, err)
		log.Println(err.Error())
		return
	}

	log.Println("trace", upclient.TraceURL(span))
}

func newUptraceClient() *uptrace.Client {
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",
	})

	return upclient
}

func redisCommands(ctx context.Context, rdb *redis.Client) error {
	if err := rdb.Set(ctx, "foo", "bar", 0).Err(); err != nil {
		return err
	}

	if err := rdb.Get(ctx, "foo").Err(); err != nil {
		return err
	}

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, "foo", "bar2", 0)
		pipe.Get(ctx, "foo")
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
