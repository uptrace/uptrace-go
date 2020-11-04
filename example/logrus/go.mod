module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/sirupsen/logrus v1.7.0
	github.com/uptrace/uptrace-go v0.4.2
	github.com/uptrace/uptrace-go/extra/otellogrus v0.1.0
	go.opentelemetry.io/otel v0.13.0
)
