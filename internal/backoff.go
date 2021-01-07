package internal

import (
	"context"
	"time"
)

func Backoff(ctx context.Context, attempt int, min, max time.Duration) error {
	if attempt == 0 {
		return nil
	}
	d := BackoffDuration(attempt, min, max)
	return Sleep(ctx, d)
}

func BackoffDuration(attempt int, min, max time.Duration) time.Duration {
	d := min
	if n := attempt - 1; n > 0 {
		d <<= n
	}
	if d < 0 || d > max {
		return max
	}
	return d
}

func Sleep(ctx context.Context, d time.Duration) error {
	done := ctx.Done()
	if done == nil {
		time.Sleep(d)
		return nil
	}

	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-done:
		return ctx.Err()
	}
}
