package uptrace

import (
	"context"
	"time"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/view"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/uptrace/uptrace-go/internal"
)

func configureMetrics(ctx context.Context, client *client, cfg *config) {
	exp, err := otlpmetricClient(ctx, client.dsn)
	if err != nil {
		return
	}

	reader := metric.NewPeriodicReader(
		exp,
		metric.WithInterval(20*time.Second),
		metric.WithTemporalitySelector(statelessTemporalitySelector),
	)
	provider := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(cfg.newResource()),
	)

	global.SetMeterProvider(provider)
	client.mp = provider

	if err := runtimemetrics.Start(); err != nil {
		internal.Logger.Printf("runtimemetrics.Start failed: %s", err)
	}
}

func otlpmetricClient(ctx context.Context, dsn *DSN) (metric.Exporter, error) {
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(dsn.OTLPHost()),
		otlpmetricgrpc.WithHeaders(map[string]string{
			// Set the Uptrace DSN here or use UPTRACE_DSN env var.
			"uptrace-dsn": dsn.String(),
		}),
		otlpmetricgrpc.WithCompressor(gzip.Name),
	}

	if dsn.Scheme == "https" {
		// Create credentials using system certificates.
		creds := credentials.NewClientTLSFromCert(nil, "")
		options = append(options, otlpmetricgrpc.WithTLSCredentials(creds))
	} else {
		options = append(options, otlpmetricgrpc.WithInsecure())
	}

	return otlpmetricgrpc.New(ctx, options...)
}

func statelessTemporalitySelector(kind view.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}
