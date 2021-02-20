package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/uptrace/uptrace-go/example/grpc/api"
	"github.com/uptrace/uptrace-go/uptrace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type helloServer struct {
	api.HelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	time.Sleep(50 * time.Millisecond)
	return &api.HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func main() {
	ctx := context.Background()

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err)
		return
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	api.RegisterHelloServiceServer(server, &helloServer{})
	if err := server.Serve(ln); err != nil {
		log.Fatal(err)
		return
	}
}
