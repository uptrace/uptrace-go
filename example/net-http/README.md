# net/http OpenTelemetry example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp)
[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-net-http/)

## Installation

To install
[otelhttp](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/net/http/otelhttp)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
```

## Example

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```

Then open http://localhost:9999
