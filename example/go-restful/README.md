# Restful OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful)

## Quickstart

Install
[otelrestful](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/emicklei/go-restful/otelrestful):

```bash
go get go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful
```

Then install OpenTelemetry filter:

```go
filter := otelrestful.OTelFilter("service-name")
restful.DefaultContainer.Filter(filter)
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

Then open http://localhost:9999
