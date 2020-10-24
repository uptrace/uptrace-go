module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20201024154641-5913df4d474e
	github.com/golang/snappy v0.0.2 // indirect
	github.com/uptrace/uptrace-go v0.3.0
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.13.0
	go.opentelemetry.io/otel v0.13.0
)
