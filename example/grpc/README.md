# gRPC OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc)

## Quickstart

Install
[otelgrpc](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/google.golang.org/grpc/otelgrpc)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc
```

To instrument gRPC client:

```go
conn, err := grpc.Dial(target,
	grpc.WithInsecure(),
	grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
)
```

To instrument gRPC server:

```go
server := grpc.NewServer(
	grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
	grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
)

```

## Example

Start the gRPC server:

```shell
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run server/server.go
```

Then start the client:

```shell
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run client/client.go
```

The server output should look like this:

```shell
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run server.go
serving on :9999
trace https://uptrace.dev/search/<project_id>?q=<trace_id>
```
