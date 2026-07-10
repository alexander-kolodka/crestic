package runpipelines

import (
	"context"
	"errors"

	"github.com/alexander-kolodka/crestic/internal/cases/backup"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

type Command struct {
	Pipelines []entity.Pipeline
	DryRun    bool
}

type Handler struct {
	runJobs handler.Handler[*backup.Command]
}

func NewHandler(runJobs handler.Handler[*backup.Command]) *Handler {
	return &Handler{
		runJobs: runJobs,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	var errs []error

	for _, pipeline := range cmd.Pipelines {
		log := logger.FromContext(ctx)
		log.Info().Str("pipeline", pipeline.Name).Msg("Processing pipeline")

		err := h.runJobs.Handle(ctx, &backup.Command{
			Jobs:   pipeline.Jobs,
			DryRun: cmd.DryRun,
		})
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
