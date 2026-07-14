package runpipelines

import (
	"context"
	"errors"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/jobs"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

type Command struct {
	Pipelines []entity.Pipeline
	DryRun    bool
}

type Handler struct {
	jobs *jobs.Runner
}

func NewHandler(jobRunner *jobs.Runner) *Handler {
	return &Handler{
		jobs: jobRunner,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	var errs []error

	if cmd.DryRun {
		ctx = jobs.WithDryRun(ctx)
	}

	for _, pipeline := range cmd.Pipelines {
		log := logger.FromContext(ctx)
		log.Info().Str("pipeline", pipeline.Name).Msg("Processing pipeline")

		err := h.jobs.Run(ctx, pipeline.Jobs)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
