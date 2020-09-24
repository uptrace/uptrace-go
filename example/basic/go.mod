module uptrace-basic-example

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.1.7
	go.opentelemetry.io/otel v0.11.0
)
