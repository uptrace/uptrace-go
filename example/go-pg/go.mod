module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/extra/pgotel/v10 v10.10.3
	github.com/go-pg/pg/v10 v10.10.3
	github.com/uptrace/uptrace-go v1.0.0-RC3
	go.opentelemetry.io/otel v1.0.0-RC2
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
)
