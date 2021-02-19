module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.2.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/uptrace/uptrace-go v0.7.7
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.17.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/trace v0.17.0
	golang.org/x/crypto v0.0.0-20210218145215-b8e89b74b9df // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210218155724-8ebf48af031b // indirect
	golang.org/x/text v0.3.5 // indirect
)
