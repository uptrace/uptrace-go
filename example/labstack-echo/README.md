# Echo instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho)

## Quickstart

To install
[otelecho](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/labstack/echo/otelecho)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho
```

Then install OpenTelemetry middleware:

```go
r := echo.New()
r.Use(otelecho.Middleware("service-name"))
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

Then open http://localhost:9999
