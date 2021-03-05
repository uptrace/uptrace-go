module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/sirupsen/logrus v1.8.0
	github.com/uptrace/uptrace-go v0.8.3
	github.com/uptrace/uptrace-go/extra/otellogrus v0.8.3
	go.opentelemetry.io/otel v0.18.0
	golang.org/x/sys v0.0.0-20210305034016-7844c3c200c3 // indirect
)
