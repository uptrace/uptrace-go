module github.com/uptrace/uptrace-go/example/go-redis-cluster

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/extra/redisotel v0.3.0
	github.com/go-redis/redis/v8 v8.7.1
	github.com/uptrace/uptrace-go v0.9.0
	go.opentelemetry.io/otel v0.19.0
)
