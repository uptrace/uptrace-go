module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/magefile/mage v1.11.0 // indirect
	github.com/sirupsen/logrus v1.8.0
	github.com/uptrace/uptrace-go v0.7.7
	github.com/uptrace/uptrace-go/extra/otellogrus v0.7.7
	go.opentelemetry.io/otel v0.17.0
	golang.org/x/sys v0.0.0-20210218155724-8ebf48af031b // indirect
)
