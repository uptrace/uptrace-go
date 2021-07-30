package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/uptrace/uptrace-go/uptrace"
)

var tracer = otel.Tracer("app_or_package_name")

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

	// MySQL query span.
	_, span := tracer.Start(ctx, "selectArticleByID")
	defer span.End()

	// Required attributes.
	span.SetAttributes(
		semconv.DBSystemMySQL,
		semconv.DBStatementKey.String("SELECT * FROM articles WHERE id = 123"),
	)

	// Optional attributes.
	span.SetAttributes(
		// This query returned 1 row.
		attribute.Int("db.rows_affected", 1),

		semconv.DBConnectionStringKey.String("localhost:3306"),
		semconv.DBUserKey.String("mysql_user"),
		semconv.DBNameKey.String("mysql_db"),

		attribute.String("code.function", "selectArticleByID"),
		attribute.String("code.filepath", "/var/lib/site/article/article.go"),
		attribute.Int("code.lineno", 33),
	)

	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
