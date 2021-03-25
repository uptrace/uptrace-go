# mongo-driver instrumentation example

[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-mongo-driver/)

To run this example you need a MongoDB server. You can start one with Docker:

```bash
make up
```

Then run the example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```
