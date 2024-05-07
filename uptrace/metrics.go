package uptrace

import (
	"context"
	"time"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/uptrace/uptrace-go/internal"
)

func configureMetrics(ctx context.Context, client *client, conf *config) {
	exp, err := otlpmetricClient(ctx, conf, client.dsn)
	if err != nil {
		internal.Logger.Printf("otlpmetricClient failed: %s", err)
		return
	}

	reader := sdkmetric.NewPeriodicReader(
		exp,
		sdkmetric.WithInterval(15*time.Second),
	)

	providerOptions := append(conf.metricOptions,
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(conf.newResource()),
	)
	provider := sdkmetric.NewMeterProvider(providerOptions...)

	otel.SetMeterProvider(provider)
	client.mp = provider

	if err := runtimemetrics.Start(); err != nil {
		internal.Logger.Printf("runtimemetrics.Start failed: %s", err)
	}
}

func otlpmetricClient(ctx context.Context, conf *config, dsn *DSN) (sdkmetric.Exporter, error) {
	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(dsn.OTLPHttpEndpoint()),
		otlpmetrichttp.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": dsn.String(),
		}),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
		otlpmetrichttp.WithTemporalitySelector(preferDeltaTemporalitySelector),
	}

	if conf.tlsConf != nil {
		options = append(options, otlpmetrichttp.WithTLSClientConfig(conf.tlsConf))
	} else if dsn.Scheme == "http" {
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	return otlpmetrichttp.New(ctx, options...)
}

func preferDeltaTemporalitySelector(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	switch kind {
	case sdkmetric.InstrumentKindCounter,
		sdkmetric.InstrumentKindObservableCounter,
		sdkmetric.InstrumentKindHistogram:
		return metricdata.DeltaTemporality
	default:
		return metricdata.CumulativeTemporality
	}
}
