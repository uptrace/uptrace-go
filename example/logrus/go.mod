module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.1.3
	github.com/uptrace/uptrace-go v1.1.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/proto/otlp v0.10.0 // indirect
	google.golang.org/genproto v0.0.0-20211101144312-62acf1d99145 // indirect
)
