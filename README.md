# Uptrace for Go

![build workflow](https://github.com/uptrace/uptrace-go/actions/workflows/build.yml/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace-go/uptrace-go)](https://pkg.go.dev/github.com/uptrace/uptrace-go/uptrace)
[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://uptrace.dev/get/opentelemetry-go)
[![Chat](https://img.shields.io/badge/-telegram-red?color=white&logo=telegram&logoColor=black)](https://t.me/uptrace)

<a href="https://uptrace.dev/get/opentelemetry-go">
  <img src="https://uptrace.dev/devicon/go-original.svg" height="200px" />
</a>

## Introduction

uptrace-go is an OpenTelemery Go distribution configured to export
[traces](https://uptrace.dev/opentelemetry/distributed-tracing) and
[metrics](https://uptrace.dev/opentelemetry/metrics) to Uptrace.

## Quickstart

Install uptrace-go:

```bash
go get github.com/uptrace/uptrace-go
```

Run the [basic example](example/basic) below using the DSN from the Uptrace project settings page.

```go
package main

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/uptrace/uptrace-go/uptrace"
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

	// Create a tracer. Usually, tracer is a global variable.
	tracer := otel.Tracer("app_or_package_name")

	// Create a root span (a trace) to measure some operation.
	ctx, main := tracer.Start(ctx, "main-operation")
	// End the span when the operation we are measuring is done.
	defer main.End()

	// The passed ctx carries the parent span (main).
	// That is how OpenTelemetry manages span relations.
	_, child1 := tracer.Start(ctx, "GET /posts/:id")
	child1.SetAttributes(
		attribute.String("http.method", "GET"),
		attribute.String("http.route", "/posts/:id"),
		attribute.String("http.url", "http://localhost:8080/posts/123"),
		attribute.Int("http.status_code", 200),
	)
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "SELECT")
	child2.SetAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("db.statement", "SELECT * FROM posts LIMIT 100"),
	)
	child2.End()

	fmt.Printf("trace: %s\n", uptrace.TraceURL(main))
}
```

## Links

- [Examples](example)
- [Documentation](https://uptrace.dev/get/opentelemetry-go)
- [OpenTelemetry Gin](https://uptrace.dev/guides/opentelemetry-gin)
- [OpenTelemetry Golang GRPC](https://uptrace.dev/guides/opentelemetry-go-grpc)
