[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/uptrace-go/extra/otellogrus)](https://pkg.go.dev/github.com/uptrace/uptrace-go/extra/otellogrus)
[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-logrus/)

# OpenTelemetry instrumentation for logrus logging

This instrumentation records logrus log messages as events on the existing span that is passed via a
`context.Context`. It does not record anything if a context does not contain a span.

## Installation

```shell
go get github.com/uptrace/uptrace-go/extra/otellogrus
```

## Usage

```go
import (
    "github.com/uptrace/uptrace-go/extra/otellogrus"
    "github.com/sirupsen/logrus"
)

// Instrument logrus.
logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
)))

// Use ctx to pass the active span.
logrus.WithContext(ctx).
	WithError(errors.New("hello world")).
	WithField("foo", "bar").
	Error("something failed")
```

See [example](/example/) for details.
