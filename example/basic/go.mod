module github.com/uptrace/uptrace-go/example/basic

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.2
	go.opentelemetry.io/otel v1.0.0
)
