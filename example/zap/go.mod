module github.com/uptrace/uptrace-go/example/zap

go 1.15

replace go.uber.org/zap => github.com/uptrace/zap v1.16.1-0.20210206140206-cdb6ad27a440

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otelzap => ../../extra/otelzap

require (
	github.com/uptrace/uptrace-go v1.0.3
	github.com/uptrace/uptrace-go/extra/otelzap v1.0.3
	go.opentelemetry.io/otel v1.0.1
	go.uber.org/zap v1.19.1
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/tools v0.1.5 // indirect
)
