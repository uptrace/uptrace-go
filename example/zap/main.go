package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/uptrace/uptrace-go/extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
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

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")

	// Use Ctx to propagate the active span.
	Logger(ctx).Error("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))

	span.End()

	fmt.Printf("\ntrace: %s\n", uptrace.TraceURL(span))
}

var (
	once   sync.Once
	logger *otelzap.Logger
)

// Logger ensures that the caller does not forget to pass the context.
func Logger(ctx context.Context) otelzap.LoggerWithCtx {
	once.Do(func() {
		l, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		logger = otelzap.New(l)
	})
	return logger.Ctx(ctx)
}
