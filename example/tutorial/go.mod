module github.com/uptrace/uptrace-go/example/tutorial

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.8.3
	go.opentelemetry.io/otel v0.18.0
)
