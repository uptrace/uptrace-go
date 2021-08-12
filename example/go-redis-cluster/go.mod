module github.com/uptrace/uptrace-go/example/go-redis-cluster

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/extra/rediscmd/v8 v8.10.0 // indirect
	github.com/go-redis/redis/extra/redisotel/v8 v8.10.0
	github.com/go-redis/redis/v8 v8.11.2
	github.com/uptrace/uptrace-go v0.21.1
	go.opentelemetry.io/otel v1.0.0-RC2
)
