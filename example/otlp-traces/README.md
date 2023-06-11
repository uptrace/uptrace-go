# OpenTelemetry Go Traces example for Uptrace

This example shows how to configure
[OTLP](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/otlp) to export traces
to Uptrace.

To run this example, you need to
[create an Uptrace project](https://uptrace.dev/get/get-started.html) and pass your project DSN via
`UPTRACE_DSN` env variable:

```go
UPTRACE_DSN="https://<token>@uptrace.dev/<project_id>" go run main.go
```
