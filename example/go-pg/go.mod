module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/extra/pgotel v0.3.1
	github.com/go-pg/pg/v10 v10.8.0
	github.com/uptrace/uptrace-go v0.19.1
	go.opentelemetry.io/otel v0.19.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492 // indirect
)
