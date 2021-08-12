module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20210707082121-9a3953d1826d
	github.com/uptrace/uptrace-go v1.0.0-RC3
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.22.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
)
