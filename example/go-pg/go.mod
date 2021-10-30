module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/extra/pgotel/v10 v10.10.6
	github.com/go-pg/pg/v10 v10.10.6
	github.com/uptrace/uptrace-go v1.0.5
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	go.opentelemetry.io/otel v1.1.0
)
