module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20211015133455-b225f9b53fa1
	github.com/golang/snappy v0.0.4 // indirect
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.26.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/trace v1.1.0
)
