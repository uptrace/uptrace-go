module github.com/uptrace/uptrace-go/example/metrics

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.1.0
	go.opentelemetry.io/otel v1.1.0
	go.opentelemetry.io/otel/metric v0.24.0
	go.opentelemetry.io/proto/otlp v0.10.0 // indirect
	golang.org/x/sys v0.0.0-20211031064116-611d5d643895 // indirect
	google.golang.org/genproto v0.0.0-20211101144312-62acf1d99145 // indirect
)
