package internal_test

import (
	"testing"

	"github.com/uptrace/uptrace-go/internal"

	"github.com/stretchr/testify/require"
)

func TestParseDSN(t *testing.T) {
	type Test struct {
		dsn  string
		otlp string
	}

	tests := []Test{
		{"https://key@uptrace.dev/1", "otlp.uptrace.dev:4317"},
		{"https://key@api.uptrace.dev/1", "otlp.uptrace.dev:4317"},
		{"https://key@localhost:1234/1", "localhost:1234"},
	}
	for _, test := range tests {
		dsn, err := internal.ParseDSN(test.dsn)
		require.NoError(t, err)
		require.Equal(t, test.otlp, dsn.OTLPEndpoint())
	}
}
