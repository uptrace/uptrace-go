package uptrace

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/uptrace/uptrace-go/internal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
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
	metricOptions  []metric.Option
}

func newConfig(opts []Option) *config {
	conf := &config{
		tracingEnabled: true,
		metricsEnabled: true,
	}

	if dsn, ok := os.LookupEnv("UPTRACE_DSN"); ok {
		conf.dsn = dsn
	}

	for _, opt := range opts {
		opt.apply(conf)
	}

	return conf
}

func (conf *config) newResource() *resource.Resource {
	if conf.resource != nil {
		if len(conf.resourceAttributes) > 0 {
			internal.Logger.Printf("WithResource overrides WithResourceAttributes (discarding %v)",
				conf.resourceAttributes)
		}
		if len(conf.resourceDetectors) > 0 {
			internal.Logger.Printf("WithResource overrides WithResourceDetectors (discarding %v)",
				conf.resourceDetectors)
		}
		return conf.resource
	}

	ctx := context.TODO()

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithDetectors(conf.resourceDetectors...),
		resource.WithAttributes(conf.resourceAttributes...))
	if err != nil {
		otel.Handle(err)
		return resource.Environment()
	}
	return res
}

//------------------------------------------------------------------------------

type Option interface {
	apply(conf *config)
}

type option func(conf *config)

func (fn option) apply(conf *config) {
	fn(conf)
}

// WithDSN configures a data source name that is used to connect to Uptrace, for example,
// `https://<token>@uptrace.dev/<project_id>`.
//
// The default is to use UPTRACE_DSN environment variable.
func WithDSN(dsn string) Option {
	return option(func(conf *config) {
		conf.dsn = dsn
	})
}

// WithServiceName configures `service.name` resource attribute.
func WithServiceName(serviceName string) Option {
	return option(func(conf *config) {
		attr := semconv.ServiceNameKey.String(serviceName)
		conf.resourceAttributes = append(conf.resourceAttributes, attr)
	})
}

// WithServiceVersion configures `service.version` resource attribute, for example, `1.0.0`.
func WithServiceVersion(serviceVersion string) Option {
	return option(func(conf *config) {
		attr := semconv.ServiceVersionKey.String(serviceVersion)
		conf.resourceAttributes = append(conf.resourceAttributes, attr)
	})
}

// WithDeploymentEnvironment configures `deployment.environment` resource attribute,
// for example, `production`.
func WithDeploymentEnvironment(env string) Option {
	return option(func(conf *config) {
		attr := semconv.DeploymentEnvironmentKey.String(env)
		conf.resourceAttributes = append(conf.resourceAttributes, attr)
	})
}

// WithResourceAttributes configures resource attributes that describe an entity that produces
// telemetry, for example, such attributes as host.name, service.name, etc.
//
// The default is to use `OTEL_RESOURCE_ATTRIBUTES` env var, for example,
// `OTEL_RESOURCE_ATTRIBUTES=service.name=myservice,service.version=1.0.0`.
func WithResourceAttributes(attrs ...attribute.KeyValue) Option {
	return option(func(conf *config) {
		conf.resourceAttributes = append(conf.resourceAttributes, attrs...)
	})
}

// WithResourceDetectors adds detectors to be evaluated for the configured resource.
func WithResourceDetectors(detectors ...resource.Detector) Option {
	return option(func(conf *config) {
		conf.resourceDetectors = append(conf.resourceDetectors, detectors...)
	})
}

// WithResource configures a resource that describes an entity that produces telemetry,
// for example, such attributes as host.name and service.name. All produced spans and metrics
// will have these attributes.
//
// WithResource overrides and replaces any other resource attributes.
func WithResource(resource *resource.Resource) Option {
	return option(func(conf *config) {
		conf.resource = resource
	})
}

func WithTLSConfig(tlsConf *tls.Config) Option {
	return option(func(conf *config) {
		conf.tlsConf = tlsConf
	})
}

//------------------------------------------------------------------------------

type TracingOption interface {
	Option
	tracing()
}

type tracingOption func(conf *config)

var _ TracingOption = (*tracingOption)(nil)

func (fn tracingOption) apply(conf *config) {
	fn(conf)
}

func (fn tracingOption) tracing() {}

// WithTracingEnabled can be used to enable/disable tracing.
func WithTracingEnabled(on bool) TracingOption {
	return tracingOption(func(conf *config) {
		conf.tracingEnabled = on
	})
}

// WithTracingDisabled disables tracing.
func WithTracingDisabled() TracingOption {
	return WithTracingEnabled(false)
}

// WithTracerProvider overwrites the default Uptrace tracer provider.
// You can use it to configure Uptrace distro to use OTLP exporter.
//
// When this option is used, you might need to call otel.SetTracerProvider
// to register the provider as the global trace provider.
func WithTracerProvider(provider *sdktrace.TracerProvider) TracingOption {
	return tracingOption(func(conf *config) {
		conf.tracerProvider = provider
	})
}

// WithTraceSampler configures a span sampler.
func WithTraceSampler(sampler sdktrace.Sampler) TracingOption {
	return tracingOption(func(conf *config) {
		conf.traceSampler = sampler
	})
}

// WithPropagator sets the global TextMapPropagator used by OpenTelemetry.
// The default is propagation.TraceContext and propagation.Baggage.
func WithPropagator(propagator propagation.TextMapPropagator) TracingOption {
	return tracingOption(func(conf *config) {
		conf.textMapPropagator = propagator
	})
}

// WithTextMapPropagator is an alias for WithPropagator.
func WithTextMapPropagator(propagator propagation.TextMapPropagator) TracingOption {
	return WithPropagator(propagator)
}

// WithPrettyPrintSpanExporter adds a span exproter that prints spans to stdout.
// It is useful for debugging or demonstration purposes.
func WithPrettyPrintSpanExporter() TracingOption {
	return tracingOption(func(conf *config) {
		conf.prettyPrint = true
	})
}

// WithBatchSpanProcessorOption specifies options used to created BatchSpanProcessor.
func WithBatchSpanProcessorOption(opts ...sdktrace.BatchSpanProcessorOption) TracingOption {
	return tracingOption(func(conf *config) {
		conf.bspOptions = append(conf.bspOptions, opts...)
	})
}

//------------------------------------------------------------------------------

type MetricsOption interface {
	Option
	metrics()
}

type metricsOption func(conf *config)

var _ MetricsOption = (*metricsOption)(nil)

func (fn metricsOption) apply(conf *config) {
	fn(conf)
}

func (fn metricsOption) metrics() {}

// WithMetricsEnabled can be used to enable/disable metrics.
func WithMetricsEnabled(on bool) MetricsOption {
	return metricsOption(func(conf *config) {
		conf.metricsEnabled = on
	})
}

// WithMetricsDisabled disables metrics.
func WithMetricsDisabled() MetricsOption {
	return WithMetricsEnabled(false)
}

func WithMetricOption(options ...metric.Option) MetricsOption {
	return metricsOption(func(conf *config) {
		conf.metricOptions = append(conf.metricOptions, options...)
	})
}
