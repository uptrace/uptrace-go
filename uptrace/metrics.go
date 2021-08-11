package uptrace

import (
	"context"
	"time"

	runtimemetrics "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric/global"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/uptrace/uptrace-go/internal"
)

func configureMetrics(ctx context.Context, client *client, cfg *config) {
	exportKindSelector := export.StatelessExportKindSelector()

	exp, err := otlpmetric.New(ctx, otlpmetricClient(client.dsn),
		otlpmetric.WithMetricExportKindSelector(exportKindSelector))
	if err != nil {
		internal.Logger.Printf("otlpmetric.New failed: %s", err)
		return
	}

	ctrl := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(),
			exportKindSelector,
		),
		controller.WithExporter(exp),
		controller.WithCollectPeriod(10*time.Second), // same as default
		controller.WithResource(cfg.newResource()),
	)

	if err := ctrl.Start(ctx); err != nil {
		internal.Logger.Printf("ctrl.Start failed: %s", err)
		return
	}

	if err := runtimemetrics.Start(); err != nil {
		internal.Logger.Printf("runtimemetrics.Start failed: %s", err)
	}

	global.SetMeterProvider(ctrl.MeterProvider())
	client.ctrl = ctrl
}

func otlpmetricClient(dsn *internal.DSN) otlpmetric.Client {
	endpoint := dsn.OTLPEndpoint()

	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(endpoint),
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

	return otlpmetricgrpc.NewClient(options...)
}