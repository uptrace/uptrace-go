module github.com/uptrace/uptrace-go/example/metrics

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.2.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/metric v0.26.0
)
