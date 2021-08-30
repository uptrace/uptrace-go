module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.5.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/uptrace/uptrace-go v1.0.0
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
)
