module github.com/uptrace/uptrace-go/example/gocsql

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/gocql/gocql v0.0.0-20210413161705-87a5d7a5ff74
	github.com/uptrace/uptrace-go v0.19.4
	go.opentelemetry.io/contrib/instrumentation/github.com/gocql/gocql/otelgocql v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/trace v0.20.0
)
