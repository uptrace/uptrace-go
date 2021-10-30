module github.com/uptrace/uptrace-go/extra/otellogrus/example

go 1.17

replace github.com/uptrace/uptrace-go/extra/otellogrus => ./..

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/uptrace-go/extra/otellogrus v1.0.5
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.1.0
	go.opentelemetry.io/otel/sdk v1.1.0
)

require (
	go.opentelemetry.io/otel/trace v1.1.0 // indirect
	golang.org/x/sys v0.0.0-20211029165221-6e7872819dc8 // indirect
)
