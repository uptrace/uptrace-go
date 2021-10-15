package main

import (
	"context"
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	_ "modernc.org/sqlite"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	db, err := otelsql.Open("sqlite", "file::memory:?cache=shared",
		otelsql.WithAttributes(semconv.DBSystemSqlite),
		otelsql.WithDBName("mydb"))
	if err != nil {
		panic(err)
	}

	// Create a tracer. Usually, tracer is a global variable.
	tracer := otel.Tracer("app_or_package_name")

	// Create a root span (a trace) to measure some operation.
	ctx, main := tracer.Start(ctx, "main-operation")
	// End the span when the operation we are measuring is done.
	defer main.End()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	var num int
	if err := db.QueryRowContext(ctx, "SELECT 42").Scan(&num); err != nil {
		panic(err)
	}

	fmt.Printf("trace: %s\n", uptrace.TraceURL(main))
}
