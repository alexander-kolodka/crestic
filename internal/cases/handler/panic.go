package handler

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/pkg/mw"
	"github.com/alexander-kolodka/crestic/internal/pkg/panix"
)

// WithPanicRecovery wraps a handler with top-level panic recovery.
// Converts panics to PanicError with full stacktrace, preventing entire application crash.
// Applied to main command handler for protection against unexpected panics.
func WithPanicRecovery[CMD any]() mw.Middleware[CMD] {
	return func(fn mw.Func[CMD]) mw.Func[CMD] {
		return func(ctx context.Context, cmd CMD) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = panix.NewPanicError(r)
				}
			}()

			return fn(ctx, cmd)
		}
	}
}
