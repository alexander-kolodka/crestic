package hooks

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/shell"
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
func (r *Runner) Execute(ctx context.Context, hooks []string) error {
	ctx = logger.WithSource(ctx, "hooks")
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
