module github.com/uptrace/uptrace-go/extra/otelzap

go 1.15

replace go.uber.org/zap => github.com/uptrace/zap v1.16.1-0.20210206140206-cdb6ad27a440

require (
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/oteltest v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0
)
