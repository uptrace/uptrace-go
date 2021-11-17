module github.com/uptrace/uptrace-go/extra/otellogrus/example

go 1.17

replace github.com/uptrace/uptrace-go/extra/otellogrus => ./..

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/uptrace-go/extra/otellogrus v1.1.0
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.2.0
	go.opentelemetry.io/otel/sdk v1.2.0
)

require (
	go.opentelemetry.io/otel/trace v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20211116061358-0a5406a5449c // indirect
)
