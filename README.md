# Uptrace Go exporter for OpenTelemetry

[![Build Status](https://travis-ci.org/uptrace/uptrace-go.svg?branch=master)](https://travis-ci.org/uptrace/uptrace-go)
[![GoDoc](https://godoc.org/github.com/uptrace/uptrace-go?status.svg)](https://pkg.go.dev/github.com/uptrace/uptrace-go?tab=doc)

## Introduction

uptrace-go is an exporter for [OpenTelemetry](https://opentelemetry.io/) that
sends your traces/spans and metrics to [Uptrace.dev](https://uptrace.dev).
Briefly the process is following:

- OpenTelemetry API is used to instrument your application with spans and
  metrics.
- OpenTelemetry SDK and this exporter are used together to export collected
  spans to Uptrace.dev.
- Uptrace.dev uses that information to help you pinpoint failures and find
  performance bottlenecks.

## Code instrumentation

You instrument your application by wrapping potentially interesting operations
with spans. Each span has:

- an operation name
- a start time and end time
- a set of key/value attributes containing data about the operation
- a set of timed events representing events, errors, logs, etc.

You create spans using a tracer:

```go
import "go.opentelemetry.io/otel/api/global"

// Create a named tracer using your repo as an identifier.
tracer := global.Tracer("github.com/username/app-name")
```

To create a span:

```go
// context is used to store and pass around active span.
ctx := context.TODO()

// Create a span and save it in the context.
ctx, span := tracer.Start(ctx, "operation-name")

// Pass the context returned from the tracer so the active span is not lost.
err := doSomeWork(ctx)

// End span when operation is completed.
span.End()
```

Alternatively you can use `WithSpan` which does roughly the same:

```go
tracer.WithSpan(ctx, "operation-name", func(ctx context.Context) error {
    return doSomeWork(ctx)
})
```

To get existing span from the context:

```go
span := trace.SpanFromContext(ctx)
```

Once you have a span you can start adding attributes:

```go
import "go.opentelemetry.io/otel/api/kv"

span.SetAttributes(
    kv.String("enduser.id", "123"),
    kv.String("enduser.role", "admin"),
)
```

or events:

```go
span.AddEvent(ctx, "log",
    kv.String("log.severity", "error"),
    kv.String("log.message", "User not found"),
    kv.String("enduser.id", "123"),
)
```

To record an error use `RecordError` which internally uses `AddEvent`:

```go
if err != nil {
    span.RecordError(ctx, err)
}
```

## Span exporter

Span exporter exports spans to Update.dev backend. It can be configured with the
following code:

```go
import (
    "github.com/uptrace/uptrace-go/uptrace"
    "go.opentelemetry.io/otel/api/global"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Resource that describes this service.
resource := resource.New(
    standard.ServiceNameKey.String("my-service"),
)

// Create a trace provider using Uptrace.dev exporter.
provider, err := sdktrace.NewProvider(
    sdktrace.WithConfig(sdktrace.Config{
        Resource:       resource,
        DefaultSampler: sdktrace.AlwaysSample(),
    }),
    uptrace.WithBatcher(&uptrace.Config{
        DSN: "", // copy your project DSN here or use UPTRACE_DSN env var
    }),
)
if err != nil {
    return err
}

// Register the provider in the system.
global.SetTraceProvider(provider)
```

## Metric exporter

Metric exporter exports metrics to Uptrace.dev backend. It can be configured
with the following code:

```go
import "github.com/uptrace/uptrace-go/upmetric"

// Resource that describes this service.
resource := resource.New(
    standard.ServiceNameKey.String("my-service"),
)

ctrl := upmetric.InstallNewPipeline(&upmetric.Config{
    DSN: "", // copy your project DSN here or use UPTRACE_DSN env var
}, push.WithResource(resource))
```
