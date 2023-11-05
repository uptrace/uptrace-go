package uptrace

import (
	"context"
	"time"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"

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
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(dsn.OTLPEndpoint()),
		otlpmetricgrpc.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": dsn.String(),
		}),
		otlpmetricgrpc.WithCompressor(gzip.Name),
		otlpmetricgrpc.WithTemporalitySelector(preferDeltaTemporalitySelector),
	}

	if conf.tlsConf != nil {
		creds := credentials.NewTLS(conf.tlsConf)
		options = append(options, otlpmetricgrpc.WithTLSCredentials(creds))
	} else if dsn.Scheme == "https" {
		// Create credentials using system certificates.
		creds := credentials.NewClientTLSFromCert(nil, "")
		options = append(options, otlpmetricgrpc.WithTLSCredentials(creds))
	} else {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	return otlpmetricgrpc.New(ctx, options...)
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
