# OpenTelemetry Go Metrics example for Uptrace

This example shows how to configure
[OTLP](https://github.com/open-telemetry/opentelemetry-go/tree/main/exporters/otlp) to export
metrics to Uptrace.

To run this example, you need to
[create an Uptrace project](https://uptrace.dev/get/get-started.html) and pass your project DSN via
`UPTRACE_DSN` env variable:

```go
UPTRACE_DSN="https://<token>@uptrace.dev/<project_id>" go run main.go
```

To view metrics, open [app.uptrace.dev](https://app.uptrace.dev/) and navigate to the Metrics ->
Explore tab.
