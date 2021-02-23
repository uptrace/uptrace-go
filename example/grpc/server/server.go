package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/uptrace/uptrace-go/example/grpc/api"
	"github.com/uptrace/uptrace-go/uptrace"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const addr = ":9999"

var upclient *uptrace.Client

type helloServer struct {
	api.HelloServiceServer
}

func (s *helloServer) SayHello(ctx context.Context, in *api.HelloRequest) (*api.HelloResponse, error) {
	log.Println("trace", upclient.TraceURL(trace.SpanFromContext(ctx)))

	time.Sleep(50 * time.Millisecond)
	return &api.HelloResponse{Reply: "Hello " + in.Greeting}, nil
}

func main() {
	ctx := context.Background()

	upclient = uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		ServiceName:    "myservice",
		ServiceVersion: "1.0.0",
	})
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	log.Println("serving on", addr)

	ln, err := net.Listen("tcp", addr)
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
