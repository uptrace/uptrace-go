module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/uptrace-go v0.19.4
	github.com/uptrace/uptrace-go/extra/otellogrus v0.19.4
	go.opentelemetry.io/otel v0.20.0
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
)
