module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.2.1
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/uptrace/uptrace-go v0.19.1
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.19.0
	go.opentelemetry.io/otel v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210324051636-2c4c8ecb7826 // indirect
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492 // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
)
