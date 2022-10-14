module github.com/uptrace/uptrace-go/example/metrics

go 1.18

replace github.com/uptrace/uptrace-go => ../..

require (
	go.opentelemetry.io/otel v1.11.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v0.32.3
	go.opentelemetry.io/otel/metric v0.32.3
	go.opentelemetry.io/otel/sdk/metric v0.32.3
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	go.opentelemetry.io/otel/sdk v1.11.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.0 // indirect
	golang.org/x/sys v0.0.0-20221010170243-090e33056c14 // indirect
)
