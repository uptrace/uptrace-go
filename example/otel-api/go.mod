module github.com/uptrace/uptrace-go/example/otel-api

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.7.7
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/trace v0.17.0
)
