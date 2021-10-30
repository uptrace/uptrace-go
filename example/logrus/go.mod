module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.1.2
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/otel v1.1.0
)
