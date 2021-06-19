module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20210515062232-b7ef815b4556
	github.com/uptrace/uptrace-go v0.20.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
)
