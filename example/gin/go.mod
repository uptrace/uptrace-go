module github.com/uptrace/uptrace-go/example/gin

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go v1.1.10 // indirect
	github.com/uptrace/uptrace-go v0.3.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.13.0
	go.opentelemetry.io/otel v0.13.0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
