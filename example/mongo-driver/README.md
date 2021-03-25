# mongo-driver OpenTelemetry instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo)

## Quickstart

Install
[otelmongo](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo)
instrumentation:

```bash
go get go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo
```

Then add OpenTelemetry monitor:

```go
opt := options.Client()
opt.Monitor = otelmongo.NewMonitor("service-name")
```

## Example

To run this example you need a MongoDB server. You can start one with Docker:

```bash
make up
```

Then run the example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```
