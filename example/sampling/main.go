package main

import (
	"context"
	"log"
	"os"

	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel/api/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	ctx := context.Background()
	upclient := setupUptrace()

	defer upclient.Close()
	defer upclient.ReportPanic(ctx)

	// Create a tracer.
	tracer := global.Tracer("github.com/your/repo")

	{
		ctx, span := tracer.Start(ctx, "trace1")
		span.End()

		_, span = tracer.Start(ctx, "child1")
		span.End()
	}

	{
		ctx, span := tracer.Start(ctx, "trace2")
		span.End()

		_, span = tracer.Start(ctx, "child1")
		span.End()
	}
}

func setupUptrace() *uptrace.Client {
	if os.Getenv("UPTRACE_DSN") == "" {
		panic("UPTRACE_DSN is empty or missing")
	}

	upclient := uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",

		Sampler: CustomSampler{Fallback: sdktrace.AlwaysSample()},

		// Pretty print spans to stdout. For debugging purposes.
		PrettyPrint: true,
	})

	return upclient
}

// CustomSampler drops some traces based on their name and uses fallback sampler otherwise.
type CustomSampler struct {
	Fallback sdktrace.Sampler
}

func (s CustomSampler) ShouldSample(params sdktrace.SamplingParameters) sdktrace.SamplingResult {
	log.Println("ShouldSample", params.Name)

	if params.Name == "trace2" {
		// Drop traces with such name.
		return sdktrace.SamplingResult{
			Decision: sdktrace.Drop,
		}
	}

	// For the rest use fallback sampler.
	return s.Fallback.ShouldSample(params)
}

func (s CustomSampler) Description() string {
	return s.Fallback.Description()
}
