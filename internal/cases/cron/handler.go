package cron

import (
	"context"

	enginecron "github.com/alexander-kolodka/crestic/internal/engine/cron"
	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/engine/pipelines"
	"github.com/alexander-kolodka/crestic/internal/entity"
)

type Command struct {
	Pipelines   []entity.Pipeline
	StateFile   string
	DryRun      bool
	Healthcheck bool
}

type Handler struct {
	pipelines *pipelines.Runner
}

// NewHandler creates a cron command Handler.
func NewHandler(pipelinesRunner *pipelines.Runner) *Handler {
	return &Handler{
		pipelines: pipelinesRunner,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	due, err := enginecron.FilterPipelinesByCron(ctx, cmd.Pipelines, cmd.StateFile)
	if err != nil {
		return err
	}

	if len(due) == 0 {
		return nil
	}

	if cmd.DryRun {
		ctx = jobs.WithDryRun(ctx)
	}

	return h.pipelines.Run(ctx, due, cmd.Healthcheck)
}
