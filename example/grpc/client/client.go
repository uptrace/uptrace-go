package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/uptrace/uptrace-go/example/grpc/api"
	"github.com/uptrace/uptrace-go/example/grpc/config"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	ctx := context.Background()

	upclient := config.SetupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from grpc client!"))

	conn, err := grpc.Dial("grpc-server:9999",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		upclient.ReportError(ctx, err)
		log.Print(err)
		return
	}
	defer func() { _ = conn.Close() }()

	client := api.NewHelloServiceClient(conn)
	if err := sayHello(ctx, client); err != nil {
		upclient.ReportError(ctx, err)
		log.Print(err)
		return
	}
}

func sayHello(ctx context.Context, client api.HelloServiceClient) error {
	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client",
		"user-id", "test-user",
	)

	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := client.SayHello(ctx, &api.HelloRequest{Greeting: "World"})
	if err != nil {
		return err
	}

	return nil
}
