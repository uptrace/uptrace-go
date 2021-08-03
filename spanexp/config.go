package spanexp

import (
	"net/http"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/uptrace/uptrace-go/internal"
)

// Config is the configuration to be used when initializing a client.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<token>@api.uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// Sampler is the default sampler used when creating new spans.
	Sampler sdktrace.Sampler

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client

	// A hook that is called before sending a span.
	BeforeSendSpan func(*Span)

	// Trace enables Uptrace exporter instrumentation.
	Trace bool

	// TraceClient enables httptrace instrumentation on the HTTP client used by Uptrace.
	TraceClient bool
}

func (cfg *Config) init() {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = internal.HTTPClient
	}

	if cfg.BeforeSendSpan == nil {
		cfg.BeforeSendSpan = func(*Span) {}
	}

	if cfg.TraceClient {
		cfg.Trace = true
	}
}
