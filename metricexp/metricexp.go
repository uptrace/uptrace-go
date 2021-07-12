package metricexp

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"github.com/uptrace/uptrace-go/internal"
)

// InstallNewPipeline instantiates a NewExportPipeline and registers it globally.
func InstallNewPipeline(
	ctx context.Context, config *Config, options ...controller.Option,
) (*controller.Controller, error) {
	ctrl, err := NewExportPipeline(config, options...)
	if err != nil {
		return nil, err
	}

	if err := ctrl.Start(ctx); err != nil {
		return nil, err
	}

	global.SetMeterProvider(ctrl.MeterProvider())

	return ctrl, nil
}

// NewExportPipeline sets up a complete export pipeline with the recommended setup.
func NewExportPipeline(
	cfg *Config, options ...controller.Option,
) (*controller.Controller, error) {
	exporter, err := NewExporter(cfg)
	if err != nil {
		return nil, err
	}

	options = append(options, controller.WithExporter(exporter))

	ctrl := controller.New(
		processor.New(
			simple.NewWithHistogramDistribution(),
			kindSelector(),
		),
		options...,
	)

	return ctrl, nil
}

//------------------------------------------------------------------------------

// Config is the configuration to be used when initializing a client.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<token>@api.uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client

	inited bool
}

func (cfg *Config) init() {
	if cfg.inited {
		return
	}
	cfg.inited = true

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = internal.HTTPClient
	}
}
