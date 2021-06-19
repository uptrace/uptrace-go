module github.com/uptrace/uptrace-go/example/labstack-echo

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/labstack/echo/v4 v4.3.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/uptrace/uptrace-go v0.20.0
	go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
)
