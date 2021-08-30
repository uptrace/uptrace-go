module github.com/uptrace/uptrace-go/example/metrics

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/metric v0.22.0
)
