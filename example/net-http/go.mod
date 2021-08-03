module github.com/uptrace/uptrace-go/example/net-http

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.21.1
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
)
