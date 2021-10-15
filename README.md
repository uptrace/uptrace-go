# Uptrace for Go

![build workflow](https://github.com/uptrace/uptrace-go/actions/workflows/build.yml/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace-go/uptrace-go)](https://pkg.go.dev/github.com/uptrace/uptrace-go/uptrace)
[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/)

<a href="https://docs.uptrace.dev/guide/go.html">
  <img src="https://docs.uptrace.dev/devicon/go-original.svg" height="200px" />
</a>

## Introduction

uptrace-go is an OpenTelemery distribution configured to export
[traces](https://docs.uptrace.dev/tracing/#spans) to Uptrace.

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
		// uptrace.WithDSN("https://<key>@api.uptrace.dev/<project_id>"),

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
	_, child1 := tracer.Start(ctx, "child1-of-main")
	child1.SetAttributes(attribute.String("key1", "value1"))
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "child2-of-main")
	child2.SetAttributes(attribute.Int("key2", 42), attribute.Float64("key3", 123.456))
	child2.End()

	fmt.Printf("trace: %s\n", uptrace.TraceURL(main))
}
```

For more details, see [documentation](https://docs.uptrace.dev/go/) and [examples](example).
