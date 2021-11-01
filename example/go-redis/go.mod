module github.com/uptrace/uptrace-go/example/go-redis

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.4
	github.com/go-redis/redis/v8 v8.11.4
	github.com/uptrace/uptrace-go v1.1.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/proto/otlp v0.10.0 // indirect
	golang.org/x/sys v0.0.0-20211031064116-611d5d643895 // indirect
)
