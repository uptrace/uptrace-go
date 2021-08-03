module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/uptrace-go v0.21.1
	github.com/uptrace/uptrace-go/extra/otellogrus v0.21.1
	go.opentelemetry.io/otel v1.0.0-RC2
)
