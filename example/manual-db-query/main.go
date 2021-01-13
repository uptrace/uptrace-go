package main

import (
	"context"
	"fmt"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/semconv"
)

var tracer = otel.Tracer("github.com/your/repo")

func main() {
	ctx := context.Background()
	upclient := newUptraceClient()

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

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
		label.Int("db.rows_affected", 1),

		semconv.DBConnectionStringKey.String("localhost:3306"),
		semconv.DBUserKey.String("mysql_user"),
		semconv.DBNameKey.String("mysql_db"),

		label.String("code.function", "selectArticleByID"),
		label.String("code.filepath", "/var/lib/site/article/article.go"),
		label.Int("code.lineno", 33),
	)

	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}

func newUptraceClient() *uptrace.Client {
	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})

	return upclient
}
