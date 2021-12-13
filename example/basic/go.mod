module github.com/uptrace/uptrace-go/example/basic

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.2.0
	go.opentelemetry.io/otel v1.3.0
)
