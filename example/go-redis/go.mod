module github.com/uptrace/uptrace-go/example/go-redis

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/v8 v8.3.0
	github.com/go-redis/redisext v0.2.1
	github.com/uptrace/uptrace-go v0.3.0
	go.opentelemetry.io/otel v0.13.0
)
