module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/extra/pgotel/v10 v10.10.5
	github.com/go-pg/pg/v10 v10.10.5
	github.com/uptrace/uptrace-go v1.0.2
	go.opentelemetry.io/otel v1.0.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
)
