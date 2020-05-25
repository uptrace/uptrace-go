# Uptrace Go exporter for OpenTelemetry

[![Build Status](https://travis-ci.org/uptrace/uptrace-go.svg?branch=master)](https://travis-ci.org/uptrace/uptrace-go)
[![GoDoc](https://godoc.org/github.com/uptrace/uptrace-go?status.svg)](https://pkg.go.dev/github.com/uptrace/uptrace-go?tab=doc)

## Introduction

uptrace-go is an exporter for [OpenTelemetry](https://opentelemetry.io/) that sends your traces and metrics to [Uptrace.dev](https://uptrace.dev).

### Trace exporter

Trace exporter can be configured with the following code:

```go
import "github.com/uptrace/uptrace-go/uptrace"

exporter := uptrace.NewExporter(&uptrace.Config{
    DSN: "", // copy your project DSN here or use UPTRACE_DSN env var
})

// Resource that describes this service.
resource := resource.New(
    standard.ServiceNameKey.String("my-service"),
)

provider, err := sdktrace.NewProvider(
    sdktrace.WithConfig(sdktrace.Config{
        Resource:       resource,
        DefaultSampler: sdktrace.AlwaysSample(),
    }),
    sdktrace.WithBatcher(exporter, sdktrace.WithMaxExportBatchSize(10000)),
)
if err != nil {
    return err
}

global.SetTraceProvider(provider)
```

## Metric exporter

Metric exporter can be configured with the following code:

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
