package upconfig

import (
	"net/http"
	"os"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config is the configuration to be used when initializing an exporter.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<key>@uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// Resource contains attributes representing an entity that produces telemetry.
	// These attributes will be copied to every span and event.
	Resource map[string]interface{}

	// Sampler is the default sampler used when creating new spans.
	Sampler sdktrace.Sampler

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client

	// Disabled disables the exporter.
	// The default is to use UPTRACE_DISABLED environment var.
	Disabled bool

	// Trace enables Uptrace exporter instrumentation.
	Trace bool

	// ClientTrace enables httptrace instrumentation on the HTTP client used by Uptrace.
	ClientTrace bool

	inited bool
}

func (cfg *Config) Init() {
	if cfg.inited {
		return
	}
	cfg.inited = true

	if _, ok := os.LookupEnv("UPTRACE_DISABLED"); ok {
		cfg.Disabled = true
		return
	}

	if cfg.DSN == "" {
		cfg.DSN = os.Getenv("UPTRACE_DSN")
	}

	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	if cfg.ClientTrace {
		cfg.Trace = true
	}
}
