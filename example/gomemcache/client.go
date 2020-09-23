package main

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/uptrace/uptrace-go/uptrace"
	otelgomemcache "go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
)

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("memcache-tracer")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	upclient = setupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from uptrace-go!"))

	mc := otelgomemcache.NewClientWithTracing(
		memcache.New("memcached-server:11211"),
	)

	ctx, s := tracer.Start(ctx, "test-operations")
	doMemcacheOperations(ctx, mc)
	s.End()
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
			"hostname": hostname,
		},
	})

	return upclient
}

func doMemcacheOperations(ctx context.Context, mc *otelgomemcache.Client) {
	mc = mc.WithContext(ctx)

	err := mc.Add(&memcache.Item{
		Key:   "foo",
		Value: []byte("bar"),
	})
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(ctx, err)
	}

	_, err = mc.Get("foo")
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(ctx, err)
	}

	_, err = mc.Get("hello")
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(ctx, err)
	}

	err = mc.Delete("foo")
	if err != nil {
		trace.SpanFromContext(ctx).RecordError(ctx, err)
	}
}
