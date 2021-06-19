module github.com/uptrace/uptrace-go/example/logrus

go 1.15

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otellogrus => ../../extra/otellogrus

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/uptrace/uptrace-go v0.20.0
	github.com/uptrace/uptrace-go/extra/otellogrus v0.20.0
	go.opentelemetry.io/otel v1.0.0-RC1
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
)
