module github.com/uptrace/uptrace-go/example/sampling

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.7.7
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/sdk v0.17.0
)
