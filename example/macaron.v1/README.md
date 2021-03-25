# Macaron instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/gopkg.in/macaron.v1/otelmacaron)

## Quickstart

To install
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

HTTP server is running at http://localhost:9999:

```bash
curl -v http://localhost:9999/profiles/admin
```
