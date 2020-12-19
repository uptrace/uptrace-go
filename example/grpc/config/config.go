package config

import (
	"github.com/uptrace/uptrace-go/uptrace"
)

func SetupUptrace() *uptrace.Client {
	return uptrace.NewClient(&uptrace.Config{
		// copy your project DSN here or use UPTRACE_DSN env var
		DSN: "",
	})
}
