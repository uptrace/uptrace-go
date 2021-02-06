package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/uptrace-go/extra/otelzap"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN:         "",
		PrettyPrint: true,
	})

	defer upclient.ReportPanic(ctx)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		otelCore := otelzap.NewOtelCore(otelzap.WithLevel(zap.NewAtomicLevelAt(zap.ErrorLevel)))
		return zapcore.NewTee(c, otelCore)
	}))

	tracer := otel.Tracer("example")
	ctx, span := tracer.Start(ctx, "main")

	// You must use Ctx to propagate the active span.
	logger.Ctx(ctx).Error("hello from zap",
		zap.Error(errors.New("hello world")),
		zap.String("foo", "bar"))

	span.End()

	// Flush the buffer and close the client.
	upclient.Close()

	fmt.Printf("trace: %s\n", upclient.TraceURL(span))
}
