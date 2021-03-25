# Gorilla Mux OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux)

## Quickstart

Install
[otelmux](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/gorilla/mux/otelmux):

```bash
go get go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux
```

Then install OpenTelemetry middleware:

```go
r := mux.NewRouter()
r.Use(otelmux.Middleware("service-name"))
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" make
```

Then open http://localhost:9999
