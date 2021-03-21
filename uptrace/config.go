package uptrace

import (
	"context"

	"github.com/uptrace/uptrace-go/internal"
	"github.com/uptrace/uptrace-go/spanexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<key>@api.uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// `service.name` resource attribute. It is merged with Config.Resource.
	// For example, `myservice`.
	ServiceName string
	// `service.version` resource attribute. It is merged with Config.Resource.
	// For example, `1.0.0`.
	ServiceVersion string
	// Any other resource attributes. They are merged with Config.Resource.
	//
	// You can also use `OTEL_RESOURCE_ATTRIBUTES` env var. For example,
	// `service.name=myservice,service.version=1.0.0`.
	ResourceAttributes []attribute.KeyValue
	// Resource contains attributes representing an entity that produces telemetry.
	// Resource attributes are copied to all spans and events.
	//
	// The default is `resource.New`.
	Resource *resource.Resource

	// Global TextMapPropagator used by OpenTelemetry.
	// The default is propagation.TraceContext and propagation.Baggage.
	TextMapPropagator propagation.TextMapPropagator

	// Sampler is the default sampler used when creating new spans.
	Sampler sdktrace.Sampler

	// A hook that is called before sending a span.
	BeforeSpanSend func(*spanexp.Span)

	// PrettyPrint pretty prints spans to the stdout.
	PrettyPrint bool

	// When specified it overwrites the default Uptrace tracer provider.
	// It can be used to configure Uptrace client to use OTLP exporter.
	TracerProvider *sdktrace.TracerProvider
}

func (cfg *Config) resource() *resource.Resource {
	res := cfg.Resource

	if res == nil {
		var err error
		res, err = resource.New(context.TODO())
		if err != nil {
			internal.Logger.Printf(context.TODO(), "resource.New failed: %s", err)
		}
	}

	return buildResource(
		res, cfg.ResourceAttributes, cfg.ServiceName, cfg.ServiceVersion)
}

func buildResource(
	res *resource.Resource,
	resourceAttributes []attribute.KeyValue,
	serviceName, serviceVersion string,
) *resource.Resource {
	var kvs []attribute.KeyValue
	kvs = append(kvs, resourceAttributes...)

	if serviceName != "" {
		kvs = append(kvs, semconv.ServiceNameKey.String(serviceName))
	}
	if serviceVersion != "" {
		kvs = append(kvs, semconv.ServiceVersionKey.String(serviceVersion))
	}

	if res == nil {
		return resource.NewWithAttributes(kvs...)
	}

	if len(kvs) > 0 {
		return resource.Merge(
			res,
			resource.NewWithAttributes(kvs...),
		)
	}

	return res
}
