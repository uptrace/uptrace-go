module github.com/uptrace/uptrace-go/example/grpc

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/golang/protobuf v1.4.3
	github.com/uptrace/uptrace-go v0.3.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.15.0
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)
