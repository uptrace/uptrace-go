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

	uptrace.ConfigureOpentelemetry(uptrace.WithDSN("https://UNKNOWN@api.uptrace.dev/2"))

	uptrace.ReportError(ctx, errors.New("hello"))
	err := uptrace.Shutdown(ctx)
	require.NoError(t, err)

	require.Equal(t,
		`send failed: status=403: project with such id and token not found (DSN="https://UNKNOWN@api.uptrace.dev/2")`,
		logger.Message())
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
