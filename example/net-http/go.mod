module github.com/uptrace/uptrace-go/example/net-http

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/uptrace/uptrace-go v0.2.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.13.0
	go.opentelemetry.io/otel v0.13.0
)
