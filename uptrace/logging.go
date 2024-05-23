package uptrace

import (
	"context"
	"time"

	"github.com/uptrace/uptrace-go/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func configureLogging(ctx context.Context, client *client, conf *config) {
	exp, err := newOtlpLogExporter(ctx, conf, client.dsn)
	if err != nil {
		internal.Logger.Printf("otlploghttp.New failed: %s", err)
		return
	}

	queueSize := queueSize()
	bspOptions := []sdklog.BatchProcessorOption{
		sdklog.WithMaxQueueSize(queueSize),
		sdklog.WithExportMaxBatchSize(queueSize),
		sdklog.WithExportInterval(10 * time.Second),
		sdklog.WithExportTimeout(10 * time.Second),
	}
	bsp := sdklog.NewBatchProcessor(exp, bspOptions...)

	var opts []sdklog.LoggerProviderOption
	opts = append(opts, sdklog.WithProcessor(bsp))
	if res := conf.newResource(); res != nil {
		opts = append(opts, sdklog.WithResource(res))
	}

	provider := sdklog.NewLoggerProvider(opts...)
	global.SetLoggerProvider(provider)
	client.lp = provider
}

func newOtlpLogExporter(
	ctx context.Context, conf *config, dsn *DSN,
) (*otlploghttp.Exporter, error) {
	options := []otlploghttp.Option{
		otlploghttp.WithEndpoint(dsn.OTLPHttpEndpoint()),
		otlploghttp.WithHeaders(map[string]string{
			"uptrace-dsn": dsn.String(),
		}),
		otlploghttp.WithCompression(otlploghttp.GzipCompression),
	}

	if conf.tlsConf != nil {
		options = append(options, otlploghttp.WithTLSClientConfig(conf.tlsConf))
	} else if dsn.Scheme == "http" {
		options = append(options, otlploghttp.WithInsecure())
	}

	return otlploghttp.New(ctx, options...)
}
