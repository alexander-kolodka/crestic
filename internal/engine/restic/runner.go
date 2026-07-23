package restic

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/engine/shell"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

type runner interface {
	Run(ctx context.Context, service string, args ...string) *shell.Result
}

type resticRunner struct {
	runner runner
}

func (r *resticRunner) Run(ctx context.Context, service string, args ...string) *shell.Result {
	ctx = logger.WithSource(ctx, "restic")
	return r.runner.Run(ctx, service, args...)
}
