# Using OTLP exporter with Uptrace

This example shows how to configure
[OTLP](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/otlp) to export traces
to Uptrace.

```go
UPTRACE_DSN="https://<token>@uptrace.dev/<project_id>" go run main.go
```
