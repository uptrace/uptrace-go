module github.com/uptrace/uptrace-go/example/go-redis

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-redis/redis/extra/redisotel v0.2.0
	github.com/go-redis/redis/v8 v8.4.0
	github.com/klauspost/compress v1.11.3 // indirect
	github.com/uptrace/uptrace-go v0.4.2
	github.com/vmihailenco/msgpack/v5 v5.0.0 // indirect
	go.opentelemetry.io/otel v0.14.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68 // indirect
)
