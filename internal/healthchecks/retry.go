package healthchecks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alexander-kolodka/crestic/internal/logger"
)

type retryableError struct {
	err error
}

func (e *retryableError) Error() string {
	return e.err.Error()
}

func (e *retryableError) Unwrap() error {
	return e.err
}

func withRetry(ctx context.Context, fn func() error) error {
	var lastErr error
	maxRetries := 3

	for attempt := range maxRetries {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			backoff := backoffDuration(attempt)
			log := logger.FromContext(ctx)
			log.Warn().
				Int("attempt", attempt).
				Int("max_retries", maxRetries).
				Dur("backoff", backoff).
				Msg("Retrying healthcheck request")
			time.Sleep(backoff)
		}

		err := fn()
		if err == nil {
			return nil
		}

		var retryable *retryableError
		if !errors.As(err, &retryable) {
			return err
		}

		lastErr = err
	}

	log := logger.FromContext(ctx)
	log.Error().
		Err(lastErr).
		Int("attempts", maxRetries).
		Msg("Healthcheck failed after all retries")

	return fmt.Errorf("healthcheck failed after %d retries: %w", maxRetries, lastErr)
}

// backoffDuration returns 2^(attempt-1) seconds, capped to avoid overflow.
//
//nolint:gosec // overflow is prevented by capping shift below 63
func backoffDuration(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	const maxAttempt = 62 // 1<<63 would overflow Duration
	shift := min(attempt-1, maxAttempt)
	return time.Second * time.Duration(1<<uint(shift))
}
