module github.com/uptrace/uptrace-go/example/go-redis

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/extra/redisotel v0.2.0
	github.com/go-redis/redis/v8 v8.6.0
	github.com/uptrace/uptrace-go v0.7.7
	go.opentelemetry.io/otel v0.17.0
)
