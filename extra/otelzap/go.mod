module github.com/uptrace/uptrace-go/extra/otelzap

go 1.15

replace go.uber.org/zap => github.com/uptrace/zap v1.16.1-0.20210206140206-cdb6ad27a440

require (
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/otel v0.16.0
	go.uber.org/zap v1.16.0
)
