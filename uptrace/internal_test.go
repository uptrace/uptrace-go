package uptrace

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPrecision(t *testing.T) {
	dur := time.Duration(math.MaxUint32) * time.Duration(prec)
	require.Equal(t, "11930h27m52.95s", dur.String())
}

func TestIDGenerator(t *testing.T) {
	ctx := context.Background()
	gen := defaultIDGenerator()

	traceID1, spanID1 := gen.NewIDs(ctx)
	traceID2, spanID2 := gen.NewIDs(ctx)
	require.NotEqual(t, traceID1, traceID2)
	require.NotEqual(t, spanID1, spanID2)

	spanID3 := gen.NewSpanID(ctx, traceID1)
	require.NotEqual(t, spanID1, spanID3)
}
