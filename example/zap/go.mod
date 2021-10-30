module github.com/uptrace/uptrace-go/example/zap

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.1.2
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/otel v1.1.0
	go.uber.org/zap v1.19.1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
)
