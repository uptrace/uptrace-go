# gRPC instrumentation example

[![Documentation](https://img.shields.io/badge/uptrace-documentation-informational)](https://docs.uptrace.dev/go/opentelemetry-grpc/)

## Running with Docker

To run this example:

```shell
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" make
```

## Running locally

Start the server:

```shell
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" go run server/server.go
```

Start the client:

```shell
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" go run client/client.go
```

The server output should look like this:

```shell
UPTRACE_DSN="https://<key>@uptrace.dev/<project_id>" go run server/server.go
serving on :9999
trace https://uptrace.dev/search/<project_id>?q=<trace_id>
```

## Other

To compile proto:

```shell
protoc -I api --go_out=plugins=grpc,paths=source_relative:./api api/hello-service.proto
```
