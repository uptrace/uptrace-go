module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.1.17
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/uptrace/uptrace-go v0.4.2
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.13.0
	go.opentelemetry.io/otel v0.13.0
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
)
