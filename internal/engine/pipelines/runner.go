package pipelines

import (
	"context"
	"errors"

	"github.com/alexander-kolodka/crestic/internal/engine/hooks"
	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/pkg/mw"
)

// Runner executes pipelines sequentially, collecting errors with errors.Join.
type Runner struct {
	jobs  *jobs.Runner
	hooks *hooks.Runner
}

// NewRunner creates a pipeline Runner.
func NewRunner(jobRunner *jobs.Runner, hooksRunner *hooks.Runner) *Runner {
	return &Runner{
		jobs:  jobRunner,
		hooks: hooksRunner,
	}
}

// Run executes pipelines and joins per-pipeline errors.
func (r *Runner) Run(ctx context.Context, pipelineList []entity.Pipeline, healthcheck bool) error {
	var errs []error

	for _, pipeline := range pipelineList {
		err := r.runPipeline(ctx, pipeline, healthcheck)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *Runner) runPipeline(ctx context.Context, pipeline entity.Pipeline, healthcheck bool) error {
	log := logger.FromContext(ctx)
	log.Info().Str("pipeline", pipeline.Name).Msg("Processing pipeline")

	fn := mw.Chain(
		r.runJobs,
		newHealthcheckMw(healthcheck),
		newHookMw(r.hooks),
	)

	return fn(ctx, pipeline)
}

func (r *Runner) runJobs(ctx context.Context, pipeline entity.Pipeline) error {
	return r.jobs.Run(ctx, pipeline.Jobs)
}
