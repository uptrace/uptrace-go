module github.com/uptrace/uptrace-go/example/sampling

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v0.15.0
	go.opentelemetry.io/otel/sdk v0.15.0
)
