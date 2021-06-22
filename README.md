# Uptrace for Go

![build workflow](https://github.com/uptrace/uptrace-go/actions/workflows/build.yml/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace-go/uptrace-go)](https://pkg.go.dev/github.com/uptrace/uptrace-go)
[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/)

<a href="https://docs.uptrace.dev/go/">
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

	uptrace.ConfigureOpentelemetry(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "myservice",
		ServiceVersion: "1.0.0",
	})
	defer uptrace.Shutdown(ctx)

	tracer := otel.Tracer("app_or_package_name")
	ctx, span := tracer.Start(ctx, "main")

	_, child1 := tracer.Start(ctx, "child1")
	child1.SetAttributes(attribute.String("key1", "value1"))
	child1.RecordError(errors.New("error1"))
	child1.End()

	_, child2 := tracer.Start(ctx, "child2")
	child2.SetAttributes(attribute.Int("key2", 42), attribute.Float64("key3", 123.456))
	child2.End()

	span.End()
	fmt.Printf("trace: %s\n", uptrace.TraceURL(span))
}
```

For more details, please see [documentation](https://docs.uptrace.dev/go/) and [examples](example).
