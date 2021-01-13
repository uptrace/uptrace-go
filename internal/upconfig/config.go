package upconfig

import (
	"context"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config is the configuration to be used when initializing a client.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<key>@uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// `service.name` resource attribute.
	ServiceName string
	// `service.version` resource attribute.
	ServiceVersion string

	// Resource contains attributes representing an entity that produces telemetry.
	// These attributes are copied to all spans and events.
	Resource *resource.Resource

	// Global TextMapPropagator used by OpenTelemetry.
	// The default is propagation.TraceContext and propagation.Baggage.
	TextMapPropagator propagation.TextMapPropagator

	// Sampler is the default sampler used when creating new spans.
	Sampler sdktrace.Sampler

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client

	// Name of the tracer used by Uptrace client.
	// The default is github.com/uptrace/uptrace-go.
	TracerName string

	// PrettyPrint pretty prints spans to the stdout.
	PrettyPrint bool

	// Disabled disables the exporter.
	// The default is to use UPTRACE_DISABLED environment var.
	Disabled bool

	// Trace enables Uptrace exporter instrumentation.
	Trace bool

	// ClientTrace enables httptrace instrumentation on the HTTP client used by Uptrace.
	ClientTrace bool

	inited bool
}

func Init(cfg *Config) {
	if cfg.inited {
		return
	}
	cfg.inited = true

	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		cfg.Disabled = true
		return
	}

	if cfg.DSN == "" {
		cfg.DSN = os.Getenv("UPTRACE_DSN")
	}

	if cfg.Resource == nil {
		resource, err := resource.New(context.TODO())
		if err == nil {
			cfg.Resource = resource
		}
	}

	if cfg.TextMapPropagator == nil {
		cfg.TextMapPropagator = propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}

	if cfg.TracerName == "" {
		cfg.TracerName = "github.com/uptrace/uptrace-go"
	}

	if cfg.Sampler == nil {
		cfg.Sampler = sdktrace.AlwaysSample()
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	if cfg.ClientTrace {
		cfg.Trace = true
	}
}
