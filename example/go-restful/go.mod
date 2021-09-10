module github.com/uptrace/uptrace-go/example/go-restful

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/emicklei/go-restful/v3 v3.5.2
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/uptrace/uptrace-go v1.0.1
	go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful v0.23.0
	go.opentelemetry.io/contrib/instrumentation/runtime v0.23.0 // indirect
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC3 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC3 // indirect
	golang.org/x/net v0.0.0-20210908191846-a5e095526f91 // indirect
	golang.org/x/sys v0.0.0-20210909193231-528a39cd75f3 // indirect
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
)
