module github.com/uptrace/uptrace-go/example/grpc

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/golang/protobuf v1.4.3
	github.com/uptrace/uptrace-go v0.8.3
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.18.0
	go.opentelemetry.io/otel/trace v0.18.0
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210305034016-7844c3c200c3 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
)
