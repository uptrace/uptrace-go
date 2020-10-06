# gRPC instrumentation example

To run this example:

```bash
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" make
```

To compile proto:

```bash
protoc -I api --go_out=plugins=grpc,paths=source_relative:./api api/hello-service.proto
```
