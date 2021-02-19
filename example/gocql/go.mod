module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20210129204804-4364a4b9cfdd
	github.com/golang/snappy v0.0.2 // indirect
	github.com/uptrace/uptrace-go v0.7.7
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.17.0
	go.opentelemetry.io/otel v0.17.0
	go.opentelemetry.io/otel/trace v0.17.0
)
