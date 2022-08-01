package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/uptrace/uptrace-go/uptrace"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func lambdaHandler(ctx context.Context) func(ctx context.Context) (interface{}, error) {
	// init aws config
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	// instrument all aws clients
	otelaws.AppendMiddlewares(&cfg.APIOptions)

	s3Client := s3.NewFromConfig(cfg)
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
		),
	}

	return func(ctx context.Context) (interface{}, error) {
		result, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return nil, err
		}

		for _, bucket := range result.Buckets {
			fmt.Println(*bucket.Name + ": " + bucket.CreationDate.Format("2006-01-02 15:04:05 Monday"))
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://api.github.com/repos/open-telemetry/opentelemetry-go/releases/latest",
			nil,
		)
		if err != nil {
			return nil, err
		}

		res, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = res.Body.Close()
		}()

		var data map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return nil, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       uptrace.TraceURL(trace.SpanFromContext(ctx)),
		}, nil
	}
}

func main() {
	ctx := context.Background()

	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		// uptrace.WithDSN("https://<token>@uptrace.dev/<project_id>"),

		uptrace.WithServiceName("myservice"),
		uptrace.WithServiceVersion("1.0.0"),
	)
	defer uptrace.Shutdown(ctx)

	tp := uptrace.TracerProvider()
	lambda.Start(otellambda.InstrumentHandler(
		lambdaHandler(ctx),
		otellambda.WithFlusher(tp),
	))
}
