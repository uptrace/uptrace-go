package uptrace

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type config struct {
	dsn string

	// Common options

	resourceAttributes []attribute.KeyValue
	resourceDetectors  []resource.Detector
	resource           *resource.Resource

	tlsConf *tls.Config

	// Tracing options

	tracingEnabled    bool
	textMapPropagator propagation.TextMapPropagator
	tracerProvider    *sdktrace.TracerProvider
	traceSampler      sdktrace.Sampler
	prettyPrint       bool
	bspOptions        []sdktrace.BatchSpanProcessorOption

	// Metrics options

	metricsEnabled bool
}

func newConfig(opts []Option) *config {
	cfg := &config{
		tracingEnabled: true,
		metricsEnabled: true,
	}

	if dsn, ok := os.LookupEnv("UPTRACE_DSN"); ok {
		cfg.dsn = dsn
	}

	for _, opt := range opts {
		opt.apply(cfg)
	}

	return cfg
}

func (cfg *config) newResource() *resource.Resource {
	if cfg.resource != nil {
		if len(cfg.resourceAttributes) > 0 {
			internal.Logger.Printf("WithResource overrides WithResourceAttributes (discarding %v)",
				cfg.resourceAttributes)
		}
		if len(cfg.resourceDetectors) > 0 {
			internal.Logger.Printf("WithResource overrides WithResourceDetectors (discarding %v)",
				cfg.resourceDetectors)
		}
		return cfg.resource
	}

	ctx := context.TODO()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithDetectors(cfg.resourceDetectors...),
		resource.WithAttributes(cfg.resourceAttributes...))
	if err != nil {
		otel.Handle(err)
		return resource.Environment()
	}
	return res
}

//------------------------------------------------------------------------------

type Option interface {
	apply(cfg *config)
}

type option func(cfg *config)

func (fn option) apply(cfg *config) {
	fn(cfg)
}

// WithDSN configures a data source name that is used to connect to Uptrace, for example,
// `https://<token>@api.uptrace.dev/<project_id>`.
//
// The default is to use UPTRACE_DSN environment variable.
func WithDSN(dsn string) Option {
	return option(func(cfg *config) {
		cfg.dsn = dsn
	})
}

// WithServiceVersion configures `service.name` resource attribute.
func WithServiceName(serviceName string) Option {
	return option(func(cfg *config) {
		attr := semconv.ServiceNameKey.String(serviceName)
		cfg.resourceAttributes = append(cfg.resourceAttributes, attr)
	})
}

// WithServiceVersion configures `service.version` resource attribute, for example, `1.0.0`.
func WithServiceVersion(serviceVersion string) Option {
	return option(func(cfg *config) {
		attr := semconv.ServiceVersionKey.String(serviceVersion)
		cfg.resourceAttributes = append(cfg.resourceAttributes, attr)
	})
}

// WithDeploymentEnvironment configures `deployment.environment` resource attribute,
// for example, `production`.
func WithDeploymentEnvironment(env string) Option {
	return option(func(cfg *config) {
		attr := semconv.DeploymentEnvironmentKey.String(env)
		cfg.resourceAttributes = append(cfg.resourceAttributes, attr)
	})
}

// WithResourceAttributes configures resource attributes that describe an entity that produces
// telemetry, for example, such attributes as host.name, service.name, etc.
//
// The default is to use `OTEL_RESOURCE_ATTRIBUTES` env var, for example,
// `OTEL_RESOURCE_ATTRIBUTES=service.name=myservice,service.version=1.0.0`.
func WithResourceAttributes(attrs ...attribute.KeyValue) Option {
	return option(func(cfg *config) {
		cfg.resourceAttributes = append(cfg.resourceAttributes, attrs...)
	})
}

// WithResourceDetectors adds detectors to be evaluated for the configured resource.
func WithResourceDetectors(detectors ...resource.Detector) Option {
	return option(func(cfg *config) {
		cfg.resourceDetectors = append(cfg.resourceDetectors, detectors...)
	})
}

// WithResource configures a resource that describes an entity that produces telemetry,
// for example, such attributes as host.name and service.name. All produced spans and metrics
// will have these attributes.
//
// WithResource overrides and replaces any other resource attributes.
func WithResource(resource *resource.Resource) Option {
	return option(func(cfg *config) {
		cfg.resource = resource
	})
}

func WithTLSConfig(tlsConf *tls.Config) Option {
	return option(func(cfg *config) {
		cfg.tlsConf = tlsConf
	})
}

//------------------------------------------------------------------------------

type TracingOption interface {
	Option
	tracing()
}

type tracingOption func(cfg *config)

var _ TracingOption = (*tracingOption)(nil)

func (fn tracingOption) apply(cfg *config) {
	fn(cfg)
}

func (fn tracingOption) tracing() {}

// WithTracingEnabled can be used to enable/disable tracing.
func WithTracingEnabled(on bool) TracingOption {
	return tracingOption(func(cfg *config) {
		cfg.tracingEnabled = on
	})
}

// WithTracingDisabled disables tracing.
func WithTracingDisabled() TracingOption {
	return WithTracingEnabled(false)
}

// TracerProvider overwrites the default Uptrace tracer provider.
// You can use it to configure Uptrace distro to use OTLP exporter.
func WithTracerProvider(provider *sdktrace.TracerProvider) TracingOption {
	return tracingOption(func(cfg *config) {
		cfg.tracerProvider = provider
	})
}

// WithTraceSampler configures a span sampler.
func WithTraceSampler(sampler sdktrace.Sampler) TracingOption {
	return tracingOption(func(cfg *config) {
		cfg.traceSampler = sampler
	})
}

// WithPropagator sets the global TextMapPropagator used by OpenTelemetry.
// The default is propagation.TraceContext and propagation.Baggage.
func WithPropagator(propagator propagation.TextMapPropagator) TracingOption {
	return tracingOption(func(cfg *config) {
		cfg.textMapPropagator = propagator
	})
}

// WithTextMapPropagator is an alias for WithPropagator.
func WithTextMapPropagator(propagator propagation.TextMapPropagator) TracingOption {
	return WithPropagator(propagator)
}

// WithPrettyPrintSpanExporter adds a span exproter that prints spans to stdout.
// It is useful for debugging or demonstration purposes.
func WithPrettyPrintSpanExporter() TracingOption {
	return tracingOption(func(cfg *config) {
		cfg.prettyPrint = true
	})
}

// WithBatchSpanProcessorOption specifies options used to created BatchSpanProcessor.
func WithBatchSpanProcessorOption(opts ...sdktrace.BatchSpanProcessorOption) TracingOption {
	return tracingOption(func(cfg *config) {
		cfg.bspOptions = append(cfg.bspOptions, opts...)
	})
}

//------------------------------------------------------------------------------

type MetricsOption interface {
	Option
	metrics()
}

type metricsOption func(cfg *config)

var _ MetricsOption = (*metricsOption)(nil)

func (fn metricsOption) apply(cfg *config) {
	fn(cfg)
}

func (fn metricsOption) metrics() {}

// WithMetricsEnabled can be used to enable/disable metrics.
func WithMetricsEnabled(on bool) MetricsOption {
	return metricsOption(func(cfg *config) {
		cfg.metricsEnabled = on
	})
}

// WithMetricsDisabled disables metrics.
func WithMetricsDisabled() MetricsOption {
	return WithMetricsEnabled(false)
}
