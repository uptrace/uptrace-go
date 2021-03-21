package spanexp

import (
	"net/http"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Config is the configuration to be used when initializing a client.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<key>@api.uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// Sampler is the default sampler used when creating new spans.
	Sampler sdktrace.Sampler

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client
	// Max number of retries when sending data to Uptrace.
	// The default is 3.
	MaxRetries int

	// A hook that is called before sending a span.
	BeforeSpanSend func(*Span)

	// Trace enables Uptrace exporter instrumentation.
	Trace bool

	// ClientTrace enables httptrace instrumentation on the HTTP client used by Uptrace.
	ClientTrace bool

	inited bool
}

func (cfg *Config) init() {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	switch cfg.MaxRetries {
	case -1:
		cfg.MaxRetries = 0
	case 0:
		cfg.MaxRetries = 3
	}

	if cfg.BeforeSpanSend == nil {
		cfg.BeforeSpanSend = func(*Span) {}
	}

	if cfg.ClientTrace {
		cfg.Trace = true
	}
}
