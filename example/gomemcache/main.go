package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("app_or_package_name")

func main() {
	flag.Parse()

	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN enar
		DSN: "",
	})
	defer uptrace.Shutdown(ctx)

	mc := otelmemcache.NewClientWithTracing(
		memcache.New(":11211"),
	)

	ctx, span := tracer.Start(ctx, "test-operations")
	defer span.End()

	doMemcacheOperations(ctx, mc)

	fmt.Println("trace", uptrace.TraceURL(span))
}

func doMemcacheOperations(ctx context.Context, mc *otelmemcache.Client) {
	mc = mc.WithContext(ctx)

	err := mc.Add(&memcache.Item{
		Key:   "foo",
		Value: []byte("bar"),
	})
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}

	_, err = mc.Get("foo")
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}

	_, err = mc.Get("hello")
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}

	err = mc.Delete("foo")
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(err)
	}
}
