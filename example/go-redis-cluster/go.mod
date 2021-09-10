module github.com/uptrace/uptrace-go/example/go-redis-cluster

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/go-redis/redis/extra/redisotel/v8 v8.11.3
	github.com/go-redis/redis/v8 v8.11.3
	github.com/uptrace/uptrace-go v1.0.1
	go.opentelemetry.io/contrib/instrumentation/runtime v0.23.0 // indirect
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC3 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC3 // indirect
	golang.org/x/net v0.0.0-20210908191846-a5e095526f91 // indirect
	golang.org/x/sys v0.0.0-20210909193231-528a39cd75f3 // indirect
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
)
