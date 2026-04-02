package uptrace

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/contrib/processors/minsev"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

func configureLogging(ctx context.Context, conf *config) *sdklog.LoggerProvider {
	var opts []sdklog.LoggerProviderOption
	if res := conf.newResource(); res != nil {
		opts = append(opts, sdklog.WithResource(res))
	}

	for _, dsn := range conf.dsn {
		dsn, err := ParseDSN(dsn)
		if err != nil {
			slog.Error("ParseDSN failed", slog.Any("err", err))
			continue
		}

		exp, err := newOtlpLogExporter(ctx, conf, dsn)
		if err != nil {
			slog.Error("otlploghttp.New failed: %s", slog.Any("err", err))
			continue
		}

		queueSize := queueSize()
		bspOptions := []sdklog.BatchProcessorOption{
			sdklog.WithMaxQueueSize(queueSize),
			sdklog.WithExportMaxBatchSize(queueSize),
			sdklog.WithExportInterval(10 * time.Second),
			sdklog.WithExportTimeout(10 * time.Second),
		}
		bsp := sdklog.NewBatchProcessor(exp, bspOptions...)

		var processor sdklog.Processor = bsp
		if conf.logMinSeverity != log.SeverityUndefined {
			// Convert from log.Severity (1-24 scale) to minsev.Severity (-8 to 15 scale).
			const sevOffset = int(log.SeverityTrace1) - int(minsev.SeverityTrace1)
			sev := minsev.Severity(int(conf.logMinSeverity) - sevOffset)
			processor = minsev.NewLogProcessor(bsp, sev)
		}
		opts = append(opts, sdklog.WithProcessor(processor))
	}

	provider := sdklog.NewLoggerProvider(opts...)
	global.SetLoggerProvider(provider)

	return provider
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
