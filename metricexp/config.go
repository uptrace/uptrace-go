package metricexp

import (
	"net/http"
	"time"
)

// Config is the configuration to be used when initializing a client.
type Config struct {
	// DSN is a data source name that is used to connect to uptrace.dev.
	// Example: https://<token>@api.uptrace.dev/<project_id>
	// The default is to use UPTRACE_DSN environment var.
	DSN string

	// HTTPClient that is used to send data to Uptrace.
	HTTPClient *http.Client
	// Max number of retries when sending data to Uptrace.
	// The default is 3.
	MaxRetries int

	// Disabled disables the exporter.
	// The default is to use UPTRACE_DISABLED environment var.
	Disabled bool

	inited bool
}

func (cfg *Config) init() {
	if cfg.inited {
		return
	}
	cfg.inited = true

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
}
