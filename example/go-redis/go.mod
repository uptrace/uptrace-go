module github.com/uptrace/uptrace-go/example/go-redis

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/extra/redisotel/v8 v8.8.2
	github.com/go-redis/redis/v8 v8.8.2
	github.com/uptrace/uptrace-go v0.19.4
	go.opentelemetry.io/otel v0.20.0
)
