package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracer = otel.Tracer("app_or_package_name")

var mysqlURI = "user:password@tcp(localhost:3306)/dbname?parseTime=true"

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

	driverName, err := otelsql.Register(
		"mysql",
		semconv.DBSystemMySQL.Value.AsString(),
	)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open(driverName, mysqlURI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := selectCurrentTime(ctx, db); err != nil {
		log.Fatal(err)
	}
}

func selectCurrentTime(ctx context.Context, db *sql.DB) *time.Time {
	ctx, span := tracer.Start(ctx, "selectCurrentTime")
	defer span.End()

	rows, err := db.QueryContext(ctx, `SELECT CURRENT_TIMESTAMP`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var currentTime time.Time
	for rows.Next() {
		err = rows.Scan(&currentTime)
		if err != nil {
			return nil
		}
	}

	fmt.Printf("current time: %s\ntrace: %s\n", currentTime, uptrace.TraceURL(span))

	return nil
}
