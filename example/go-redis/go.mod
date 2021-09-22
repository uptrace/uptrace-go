module github.com/uptrace/uptrace-go/example/go-redis

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/uptrace/uptrace-go v1.0.1
	go.opentelemetry.io/otel v1.0.0
)
