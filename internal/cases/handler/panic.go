package handler

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/panix"
)

// WithPanicRecovery wraps a handler with top-level panic recovery.
// Converts panics to PanicError with full stacktrace, preventing entire application crash.
// Applied to main command handler for protection against unexpected panics.
func WithPanicRecovery[CMD any]() Middleware[CMD] {
	return func(h Handler[CMD]) Handler[CMD] {
		return NewHandler(func(ctx context.Context, cmd CMD) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = panix.NewPanicError(r)
				}
			}()

			return h.Handle(ctx, cmd)
		})
	}
}
