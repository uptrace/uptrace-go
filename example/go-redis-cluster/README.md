# go-redis cluster instrumentation example

[![PkgGoDev](https://pkg.go.dev/badge/github.com/go-redis/redis/tree/master/extra/redisotel)](https://pkg.go.dev/github.com/go-redis/redis/tree/master/extra/redisotel)

## Quickstart

Install [redisotel](https://github.com/go-redis/redis/tree/master/extra/redisotel) instrumentation:

```bash
go get github.com/go-redis/redis/extra/redisotel
```

Then add OpenTelemetry hook:

```go
rdb := redis.NewClient(&redis.Options{
    Addr: "redis-server:6379",
})
rdb.AddHook(redisext.OpenTelemetryHook{})
```

## Example

To run this example you need a Redis Cluster. You can start one with Docker:

```bash
make up
```

To run this example:

```bash
UPTRACE_DSN="https://<token>@api.uptrace.dev/<project_id>" make
```
