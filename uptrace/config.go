package uptrace

import "github.com/uptrace/uptrace-go/spanexp"

type Config = spanexp.Config

type Option func(*Config)

// WithFilter is a helper that adds the filter to a Config.
func WithFilter(filter func(*spanexp.Span) bool) Option {
	return func(cfg *Config) {
		cfg.Filters = append(cfg.Filters, filter)
	}
}
