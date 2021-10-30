module github.com/uptrace/uptrace-go/example/net-http

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.26.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
)
