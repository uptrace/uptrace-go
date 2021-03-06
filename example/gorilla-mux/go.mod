module github.com/uptrace/uptrace-go/example/gorilla-mux

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gorilla/mux v1.8.0
	github.com/uptrace/uptrace-go v0.20.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
)
