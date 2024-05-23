package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	ctx := context.Background()

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<token>@uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	// Send buffered spans and free resources.
	defer uptrace.Shutdown(ctx)

	tracer := otel.Tracer("app_or_package_name")
	logger := otelslog.NewLogger("app_or_package_name")

	ctx, main := tracer.Start(ctx, "main-operation", trace.WithSpanKind(trace.SpanKindServer))
	defer main.End()

	logger.WarnContext(ctx, "xxx yyy")
	logger.ErrorContext(ctx, "hello world", slog.String("error", "error message"))

	fmt.Printf("trace: %s\n", uptrace.TraceURL(main))
}
