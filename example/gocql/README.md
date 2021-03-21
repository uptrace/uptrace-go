# Gocql instrumentation example

[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-gocql/)

To run this example you need a Cassandra server. You can start one with Docker:

```bash
make up
```

Then run the example:

```bash
UPTRACE_DSN="https://<key>@api.uptrace.dev/<project_id>" go run main.go
```
