package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "myservice",
		ServiceVersion: "1.0.0",
	})

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main", trace.WithSpanKind(trace.SpanKindProducer))

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(100 * time.Millisecond)

		_, child1 := tracer.Start(ctx, "child1", trace.WithSpanKind(trace.SpanKindConsumer))
		child1.SetAttributes(attribute.String("key1", "value1"))
		child1.RecordError(errors.New("error1"))

		time.Sleep(500 * time.Millisecond)

		child1.End()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		time.Sleep(50 * time.Millisecond)

		_, child2 := tracer.Start(ctx, "child2", trace.WithSpanKind(trace.SpanKindConsumer))
		child2.SetAttributes(attribute.Int("key2", 42), attribute.Float64("key3", 123.456))

		time.Sleep(700 * time.Millisecond)

		child2.End()
	}()

	time.Sleep(10 * time.Millisecond)
	_, child3 := tracer.Start(ctx, "child3")
	time.Sleep(150 * time.Millisecond)
	child3.End()

	time.Sleep(100 * time.Millisecond)

	span.End()

	wg.Wait()

	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}
