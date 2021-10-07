module github.com/uptrace/uptrace-go/example/go-restful

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/emicklei/go-restful/v3 v3.7.1
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/uptrace/uptrace-go v1.0.2
	go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful v0.25.0
	go.opentelemetry.io/otel v1.0.1
)
