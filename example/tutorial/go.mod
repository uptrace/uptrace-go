module github.com/uptrace/uptrace-go/example/tutorial

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/otel v1.1.0
)
