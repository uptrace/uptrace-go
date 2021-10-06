module github.com/uptrace/opentelemetry-go/otelzap/example

go 1.17

replace github.com/uptrace/opentelemetry-go/otelzap => ./..

require (
	github.com/uptrace/opentelemetry-go/otelzap v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	go.uber.org/zap v1.19.1
)

require (
	go.opentelemetry.io/otel/trace v1.0.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
)
