module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20210303210847-f18e0979d243
	github.com/golang/snappy v0.0.3 // indirect
	github.com/uptrace/uptrace-go v0.8.3
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.18.0
	go.opentelemetry.io/otel v0.18.0
	go.opentelemetry.io/otel/trace v0.18.0
)
