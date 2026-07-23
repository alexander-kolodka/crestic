package backup

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/entity"
)

type Command struct {
	Jobs   []entity.Job
	DryRun bool
}

type Handler struct {
	jobs *jobs.Runner
}

// NewHandler creates a backup command Handler.
func NewHandler(jobRunner *jobs.Runner) *Handler {
	return &Handler{
		jobs: jobRunner,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	if cmd.DryRun {
		ctx = jobs.WithDryRun(ctx)
	}

	return h.jobs.Run(ctx, cmd.Jobs)
}
