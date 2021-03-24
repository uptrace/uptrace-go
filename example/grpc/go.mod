module github.com/uptrace/uptrace-go/example/grpc

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/golang/protobuf v1.5.1
	github.com/uptrace/uptrace-go v0.19.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.19.0
	go.opentelemetry.io/otel/trace v0.19.0
	golang.org/x/net v0.0.0-20210324051636-2c4c8ecb7826 // indirect
	golang.org/x/sys v0.0.0-20210324051608-47abb6519492 // indirect
	google.golang.org/genproto v0.0.0-20210323160006-e668133fea6a // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.26.0
)
