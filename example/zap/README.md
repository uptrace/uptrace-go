# Zap instrumentation example for Uptrace and OpenTelemetry

[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-zap/)

To run this example:

```bash
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" go run main.go
```

**Note** that this example requires patching the zap package:

```go
go mod edit -replace go.uber.org/zap=github.com/uptrace/zap@master
```
