package hooks

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

const (
	PhaseBefore  = "before"
	PhaseSuccess = "success"
	PhaseFailure = "failure"
)

// Runner executes shell hook commands sequentially.
type Runner struct {
	shell *shell.Executor
}

// New creates a hooks Runner.
func New(shellExecutor *shell.Executor) *Runner {
	return &Runner{shell: shellExecutor}
}

// Execute runs hooks sequentially via sh -c and returns on the first failure.
// Empty cmds are a no-op (no log).
func (r *Runner) Execute(ctx context.Context, phase string, hooks []string) error {
	if len(hooks) == 0 {
		return nil
	}

	ctx = logger.WithSource(ctx, "hooks")
	log := logger.FromContext(ctx)
	log.Info().Str("phase", phase).Msg("Running hooks")

	for _, hook := range hooks {
		result := r.shell.Run(ctx, "sh", "-c", hook)
		if result.Error != nil {
			return fmt.Errorf(
				`hook failed "%s" [exit code %d]: %w`,
				hook,
				result.ExitCode,
				result.Error,
			)
		}
	}
	return nil
}
