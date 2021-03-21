# Gin instrumentation example

[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-gin-gonic/)

To run this example:

```bash
UPTRACE_DSN="https://<key>@api.uptrace.dev/<project_id>" make
```

HTTP server is running at http://localhost:9999:

```bash
curl -v http://localhost:9999/profiles/admin
```
