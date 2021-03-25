# Macaron OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron)

## Quickstart

Install
[otelmacaron](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/gopkg.in/macaron.v1/otelmacaron)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron
```

Then install OpenTelemetry middleware:

```go
m := macaron.Classic()
m.Use(otelmacaron.Middleware("service-name"))
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" make
```

Then open http://localhost:9999
