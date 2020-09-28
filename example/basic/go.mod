module uptrace-basic-example

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.2.0
	go.opentelemetry.io/otel v0.12.0
)
