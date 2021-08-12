module github.com/uptrace/uptrace-go/example/grpc

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/golang/protobuf v1.5.2
	github.com/uptrace/uptrace-go v1.0.0-RC3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.22.0
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)
