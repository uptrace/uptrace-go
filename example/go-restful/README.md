# Restful instrumentation example

[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-go-restful/)

To run this example:

```bash
UPTRACE_DSN="https://<key>@api.uptrace.dev/<project_id>" go run main.go
```

HTTP server is running at http://localhost:9999:

```bash
curl -v http://localhost:9999/profiles/admin
```
