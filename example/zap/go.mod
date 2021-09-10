module github.com/uptrace/uptrace-go/example/zap

go 1.15

replace go.uber.org/zap => github.com/uptrace/zap v1.16.1-0.20210206140206-cdb6ad27a440

replace github.com/uptrace/uptrace-go => ../..

replace github.com/uptrace/uptrace-go/extra/otelzap => ../../extra/otelzap

require (
	github.com/uptrace/uptrace-go v1.0.0
	github.com/uptrace/uptrace-go/extra/otelzap v1.0.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.23.0 // indirect
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC3 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC3 // indirect
	go.uber.org/zap v1.19.1
	golang.org/x/net v0.0.0-20210908191846-a5e095526f91 // indirect
	golang.org/x/sys v0.0.0-20210909193231-528a39cd75f3 // indirect
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
)
