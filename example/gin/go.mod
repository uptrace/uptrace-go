module github.com/uptrace/uptrace-go/example/gin

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/ugorji/go v1.2.6 // indirect
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.26.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
