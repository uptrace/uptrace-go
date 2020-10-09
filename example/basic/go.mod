module uptrace-basic-example

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v0.2.0
	go.opentelemetry.io/otel v0.13.0
	golang.org/x/net v0.0.0-20200927032502-5d4f70055728 // indirect
	golang.org/x/sys v0.0.0-20200926100807-9d91bd62050c // indirect
)
