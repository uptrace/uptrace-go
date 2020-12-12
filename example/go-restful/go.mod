module github.com/uptrace/uptrace-go/example/gin

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/emicklei/go-restful/v3 v3.3.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/uptrace/uptrace-go v0.3.0
	go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful v0.15.0
	go.opentelemetry.io/otel v0.15.0
)
