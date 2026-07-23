package runpipelines

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/engine/pipelines"
	"github.com/alexander-kolodka/crestic/internal/entity"
)

type Command struct {
	Pipelines   []entity.Pipeline
	DryRun      bool
	Healthcheck bool
}

type Handler struct {
	pipelines *pipelines.Runner
}

func NewHandler(pipelinesRunner *pipelines.Runner) *Handler {
	return &Handler{
		pipelines: pipelinesRunner,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	if cmd.DryRun {
		ctx = jobs.WithDryRun(ctx)
	}

	return h.pipelines.Run(ctx, cmd.Pipelines, cmd.Healthcheck)
}
