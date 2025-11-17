package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofrs/flock"

	"github.com/alexander-kolodka/crestic/internal/logger"
)

// WithLock wraps a handler with file-based locking to prevent concurrent execution.
// The lock file is created in ~/.crestic/ directory.
// The lock is released on normal return, on panic, and when ctx is canceled (SIGINT/SIGTERM).
func WithLock[CMD any](lockFile string) Middleware[CMD] {
	return func(h Handler[CMD]) Handler[CMD] {
		return NewHandler(func(ctx context.Context, cmd CMD) error {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}

			cresticDir := filepath.Join(home, ".crestic")
			err = os.MkdirAll(cresticDir, 0o750)
			if err != nil {
				return fmt.Errorf("failed to create .crestic directory: %w", err)
			}

			lockPath := filepath.Join(cresticDir, lockFile)
			fileLock := flock.New(lockPath)

			locked, err := fileLock.TryLock()
			if err != nil {
				return fmt.Errorf("failed to acquire lock: %w", err)
			}
			if !locked {
				return fmt.Errorf("another process is already running (lock: %s)", lockPath)
			}

			// Ensure we release the lock exactly once in all code paths.
			var once sync.Once
			release := func() {
				once.Do(func() {
					uErr := fileLock.Unlock()
					if uErr != nil {
						log := logger.FromContext(ctx)
						log.Err(uErr).Msg("failed to release lock")
					}

					_ = os.Remove(lockPath)
				})
			}

			// 1) Release on panic or normal return.
			defer release()

			// 2) Release on context cancellation (Ctrl-C / SIGTERM).
			go func() {
				<-ctx.Done()
				release()
			}()

			return h.Handle(ctx, cmd)
		})
	}
}
