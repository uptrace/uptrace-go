# Gomemcache instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache)](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache)

## Quickstart

Install
[otelmemcache](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/master/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache):

```bash
go get go.opentelemetry.io/contrib/instrumentation/github.com/bradfitz/gomemcache/memcache/otelmemcache
```

Then use `NewClientWithTracing` to wrap a client:

```go
mc := otelmemcache.NewClientWithTracing(
    memcache.New("localhost:11211"),
)
```

And use `mc.WithContext` to pass around active span.

## Example

To run this example you need a memcached server. You can start one with Docker:

```
make up
```

Then run the example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" go run main.go
```
