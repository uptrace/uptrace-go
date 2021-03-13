package uptrace

import "github.com/uptrace/uptrace-go/internal"

// SetLogger sets the logger to the given one.
func SetLogger(logger internal.ILogger) {
	internal.Logger = logger
}
