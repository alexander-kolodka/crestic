package backup

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/engine/pipelines"
	"github.com/alexander-kolodka/crestic/internal/entity"
)

type Command struct {
	Pipelines   []entity.Pipeline
	Jobs        []entity.Job
	DryRun      bool
	Healthcheck bool
}

type Handler struct {
	pipelines *pipelines.Runner
	jobs      *jobs.Runner
}

// NewHandler creates a backup command Handler.
func NewHandler(pipelinesRunner *pipelines.Runner, jobRunner *jobs.Runner) *Handler {
	return &Handler{
		pipelines: pipelinesRunner,
		jobs:      jobRunner,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	if cmd.DryRun {
		ctx = jobs.WithDryRun(ctx)
	}

	if len(cmd.Pipelines) > 0 {
		return h.pipelines.Run(ctx, cmd.Pipelines, cmd.Healthcheck)
	}

	return h.jobs.Run(ctx, cmd.Jobs)
}
