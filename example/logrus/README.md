# Logrus OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/uptrace-go/extra/otellogrus)](https://pkg.go.dev/github.com/uptrace/uptrace-go/extra/otellogrus)

## Quickstart

To install [otellogrus](https://github.com/uptrace/uptrace-go/tree/master/extra/otellogrus)
instrumentation:

```bash
go get github.com/uptrace/uptrace-go/extra/otellogrus
```

Then add OpenTelemetry hook:

```go
logrus.AddHook(otellogrus.NewLoggingHook())
```

And use `WithContext` to propagate the active span:

```go
logrus.WithContext(ctx).
    WithError(errors.New("hello world")).
    WithField("foo", "bar").
    Error("hello from logrus")
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```
