module github.com/uptrace/uptrace-go/example/zap

go 1.15

replace go.uber.org/zap => github.com/uptrace/zap v1.16.1-0.20210206140206-cdb6ad27a440

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otelzap => ../../extra/otelzap

require (
	github.com/uptrace/uptrace-go v0.19.4
	github.com/uptrace/uptrace-go/extra/otelzap v0.19.4
	go.opentelemetry.io/otel v0.20.0
	go.uber.org/zap v1.16.0
)
