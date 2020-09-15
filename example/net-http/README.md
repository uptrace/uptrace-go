# net/http instrumentation example

To run this example:

```bash
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" make
```

HTTP server is running at http://localhost:9999:

```bash
curl -v http://localhost:9999/profiles/admin
```
