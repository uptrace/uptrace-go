package uptrace

import (
	"context"
	"log/slog"
	"time"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"

	"github.com/uptrace/uptrace-go/internal"
)

func configureMetrics(ctx context.Context, conf *config) *sdkmetric.MeterProvider {
	opts := conf.metricOptions
	if res := conf.newResource(); res != nil {
		opts = append(opts, sdkmetric.WithResource(res))
	}

	for _, dsn := range conf.dsn {
		dsn, err := ParseDSN(dsn)
		if err != nil {
			slog.Error("ParseDSN failed", slog.Any("err", err))
			continue
		}

		exp, err := otlpmetricClient(ctx, conf, dsn)
		if err != nil {
			internal.Logger.Printf("otlpmetricClient failed: %s", err)
			continue
		}

		reader := sdkmetric.NewPeriodicReader(
			exp,
			sdkmetric.WithInterval(15*time.Second),
		)
		opts = append(opts, sdkmetric.WithReader(reader))
	}

	provider := sdkmetric.NewMeterProvider(opts...)
	otel.SetMeterProvider(provider)

	if err := runtimemetrics.Start(); err != nil {
		slog.Error("runtimemetrics.Start failed", slog.Any("err", err))
	}

	return provider
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
