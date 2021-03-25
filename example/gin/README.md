# Gin OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin)

## Quickstart

Install
[otelgin](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/gin-gonic/gin/otelgin)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin
```

Then install the OpenTelemetry middleware:

```go
router := gin.Default()
router.Use(otelgin.Middleware("service-name"))
```

To propagate active span through the app, use `context.Context` from the `http.Request` (not
`gin.Context`). To measure HTML template rendering, use `otelgin.HTML` helper.

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

Then open http://localhost:9999
