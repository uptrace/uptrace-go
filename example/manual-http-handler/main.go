package main

import (
	"context"
	"fmt"
	"net/http"

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

	// HTTP handler span.
	_, span := tracer.Start(ctx, "GET /articles/:articleID")
	defer span.End()

	// Required attributes.
	span.SetAttributes(
		semconv.HTTPMethodKey.String("GET"),
		semconv.HTTPRouteKey.String("/articles/:articleID"),
	)

	// Optional attributes.
	span.SetAttributes(
		semconv.HTTPTargetKey.String("/articles/123"),
		semconv.HTTPStatusCodeKey.Int(http.StatusOK),
		semconv.HTTPUserAgentKey.String("Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0"),
		semconv.HTTPClientIPKey.String("127.0.0.1"),

		attribute.String("code.function", "articleEndpoint"),
		attribute.String("code.filepath", "/var/lib/site/article/article_api.go"),
		attribute.Int("code.lineno", 55),
	)

	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
