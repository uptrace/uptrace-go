module github.com/uptrace/uptrace-go/extra/otelzap/example

go 1.17

replace github.com/uptrace/uptrace-go/extra/otelzap => ./..

require (
	github.com/uptrace/uptrace-go/extra/otelzap v1.0.4
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.1
	go.opentelemetry.io/otel/sdk v1.0.1
	go.uber.org/zap v1.19.1
)

require (
	go.opentelemetry.io/otel/trace v1.0.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
)
