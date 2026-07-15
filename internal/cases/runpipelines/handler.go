package runpipelines

import (
	"context"
	"errors"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/hooks"
	"github.com/alexander-kolodka/crestic/internal/jobs"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/mw"
)

type Command struct {
	Pipelines   []entity.Pipeline
	DryRun      bool
	Healthcheck bool
}

type Handler struct {
	jobs  *jobs.Runner
	hooks *hooks.Runner
}

type Healthcheck interface {
	Start(ctx context.Context, rid, body string) error
	Success(ctx context.Context, rid, body string) error
	Fail(ctx context.Context, rid, body string) error
}

func NewHandler(jobRunner *jobs.Runner, hooksRunner *hooks.Runner) *Handler {
	return &Handler{
		jobs:  jobRunner,
		hooks: hooksRunner,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	var errs []error

	if cmd.DryRun {
		ctx = jobs.WithDryRun(ctx)
	}

	for _, pipeline := range cmd.Pipelines {
		err := h.runPipeline(ctx, cmd, pipeline)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (h *Handler) runPipeline(ctx context.Context, cmd *Command, pipeline entity.Pipeline) error {
	log := logger.FromContext(ctx)
	log.Info().Str("pipeline", pipeline.Name).Msg("Processing pipeline")

	fn := mw.Chain(
		h.runJobs,
		newHealthcheckMw(cmd),
		newHookMw(h.hooks),
	)

	return fn(ctx, pipeline)
}

func (h *Handler) runJobs(ctx context.Context, pipeline entity.Pipeline) error {
	return h.jobs.Run(ctx, pipeline.Jobs)
}
