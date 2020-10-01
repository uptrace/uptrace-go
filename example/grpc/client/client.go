package main

import (
	"context"
	"errors"
	"time"
	"log"

	"grpc/api"
	"grpc/config"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/api/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	upclient *uptrace.Client
	tracer   = global.Tracer("grpc-client-tracer")
)

func main() {
	ctx := context.Background()

	upclient = config.SetupUptrace()
	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	upclient.ReportError(ctx, errors.New("hello from grpc client!"))

	dialOption := grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor(tracer))
	conn, err := grpc.Dial("grpc-server:9999", grpc.WithInsecure(), dialOption)
	if err != nil {
		upclient.ReportError(ctx, err)
		log.Fatal(err)
		return
	}
	defer func() { _ = conn.Close() }()

	client := api.NewHelloServiceClient(conn)
	err = sayHello(ctx, client)
	if err != nil {
		upclient.ReportError(ctx, err)
		log.Fatal(err)
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