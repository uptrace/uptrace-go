[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/uptrace-go/extra/otelzap)](https://pkg.go.dev/github.com/uptrace/uptrace-go/extra/otelzap)

# OpenTelemetry Go instrumentation for Zap logging library

This instrumentation records Zap log messages as events on the existing span that is passed via a
`context.Context`. It does not record anything if a context does not contain a span.

## Installation

```shell
go get github.com/uptrace/uptrace-go/extra/otelzap
```

## Usage

You need to create a `otelzap.Logger` that wraps a `zap.Logger` and provides context-aware logging
API.

```go
import (
    "go.uber.org/zap"
    "github.com/uptrace/uptrace-go/extra/otelzap"
)

// Wrap zap logger.
log := otelzap.New(zap.L())

// And then pass ctx to propagate the span.
log.Ctx(ctx).Error("hello from zap",
	zap.Error(errors.New("hello world")),
	zap.String("foo", "bar"))

// Alternatively.
log.ErrorContext(ctx, "hello from zap",
	zap.Error(errors.New("hello world")),
	zap.String("foo", "bar"))
```

See [example](/example/) for details.

## Options

`otelzap.New` accepts a couple of
[options](https://pkg.go.dev/github.com/uptrace/uptrace-go/extra/otelzap#Option):

- `otelzap.WithMinLevel(zap.WarnLevel)` sets the minimal zap logging level on which the log message
  is recorded on the span.
- `otelzap.WithErrorStatusLevel(zap.ErrorLevel)` sets the minimal zap logging level on which the
  span status is set to codes.Error.
- `otelzap.WithCaller(true)` configures the logger to annotate each event with the filename, line
  number, and function name of the caller. Enabled by default.
- `otelzap.WithStackTrace(true)` configures the logger to capture logs with a stack trace. Disabled
  by default.
