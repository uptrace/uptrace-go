package uptrace_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/uptrace/uptrace-go/uptrace"
)

func TestInvalidDSN(t *testing.T) {
	t.Skip("overwrites global tracer provider")

	var logger Logger
	uptrace.SetLogger(&logger)

	uptrace.ConfigureOpentelemetry(uptrace.WithDSN("dsn"))

	require.Equal(t,
		`Uptrace is disabled: DSN="dsn" does not have a token`,
		logger.Message())
}

func TestUnknownToken(t *testing.T) {
	t.Skip("overwrites global tracer provider")

	ctx := context.Background()

	var logger Logger
	uptrace.SetLogger(&logger)

	uptrace.ConfigureOpentelemetry(uptrace.WithDSN("https://UNKNOWN@uptrace.dev/2"))

	uptrace.ReportError(ctx, errors.New("hello"))
	err := uptrace.Shutdown(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), `project with token="UNKNOWN" doesn't exist`)
}

//------------------------------------------------------------------------------

type Logger struct {
	msgs []string
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(format, args...))
}

func (l *Logger) Message() string {
	if len(l.msgs) == 0 {
		return ""
	}
	return l.msgs[len(l.msgs)-1]
}
