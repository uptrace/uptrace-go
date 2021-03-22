module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/uptrace-go v0.9.0
	github.com/uptrace/uptrace-go/extra/otellogrus v0.9.0
	go.opentelemetry.io/otel v0.19.0
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
)
