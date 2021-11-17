module github.com/uptrace/uptrace-go/example/basic

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.1.0
	go.opentelemetry.io/otel v1.2.0
)
