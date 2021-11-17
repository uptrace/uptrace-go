module github.com/uptrace/uptrace-go/extra/otelzap/example

go 1.17

replace github.com/uptrace/uptrace-go/extra/otelzap => ./..

require (
	github.com/uptrace/uptrace-go/extra/otelzap v1.2.0
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0
	go.opentelemetry.io/otel/sdk v1.2.0
	go.uber.org/zap v1.19.1
)

require (
	go.opentelemetry.io/otel/trace v1.2.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20211116061358-0a5406a5449c // indirect
)
