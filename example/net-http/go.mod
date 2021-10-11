module github.com/uptrace/uptrace-go/example/net-http

go 1.13

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/uptrace/uptrace-go v1.0.5
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.25.0
	go.opentelemetry.io/otel v1.0.1
	go.opentelemetry.io/otel/trace v1.0.1
	golang.org/x/net v0.0.0-20211008194852-3b03d305991f // indirect
	google.golang.org/genproto v0.0.0-20211008145708-270636b82663 // indirect
)
