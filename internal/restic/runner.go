package restic

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

type runner interface {
	Run(ctx context.Context, service string, args ...string) *shell.Result
}

type resticRunner struct {
	runner runner
}

func (r *resticRunner) Run(ctx context.Context, service string, args ...string) *shell.Result {
	ctx = logger.WithSource(ctx, "restic")

	if logger.IsJSONMode(ctx) && len(args) > 0 {
		newArgs := make([]string, 0, len(args)+1)
		newArgs = append(newArgs, args[0])
		newArgs = append(newArgs, "--json")
		newArgs = append(newArgs, args[1:]...)
		args = newArgs
	}

	return r.runner.Run(ctx, service, args...)
}
