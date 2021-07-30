module github.com/uptrace/uptrace-go/example/grpc

go 1.14

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/golang/protobuf v1.5.2
	github.com/uptrace/uptrace-go v0.20.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.21.0
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	google.golang.org/genproto v0.0.0-20210617175327-b9e0b3197ced // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1
)
