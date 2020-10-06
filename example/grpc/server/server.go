package main

import (
	"context"
	"errors"
	"grpc/api"
	"grpc/config"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
)

var tracer = global.Tracer("grpc-server-tracer")

type helloServer struct {
	api.HelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	time.Sleep(50 * time.Millisecond)
	return &api.HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func main() {
	ctx := context.Background()

	upclient := config.SetupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from grpc server!"))

	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		upclient.ReportError(ctx, err)
		log.Print(err)
		return
	}

	opt := grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(tracer))
	server := grpc.NewServer(opt)

	api.RegisterHelloServiceServer(server, &helloServer{})
	if err := server.Serve(lis); err != nil {
		upclient.ReportError(ctx, err)
		log.Print(err)
		return
	}
}
