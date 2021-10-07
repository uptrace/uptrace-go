module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.6.1
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/uptrace/uptrace-go v1.0.3
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.25.0
	go.opentelemetry.io/otel/trace v1.0.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)
