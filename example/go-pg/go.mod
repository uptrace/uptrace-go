module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/extra/pgotel v0.2.0
	github.com/go-pg/pg/v10 v10.7.7
	github.com/uptrace/uptrace-go v0.7.7
	go.opentelemetry.io/otel v0.17.0
)
