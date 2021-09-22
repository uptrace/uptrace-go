module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.6.0
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/uptrace/uptrace-go v1.0.1
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.24.0
	go.opentelemetry.io/otel/trace v1.0.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)
