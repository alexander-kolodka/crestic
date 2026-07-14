package runpipelines

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
	"github.com/alexander-kolodka/crestic/internal/jobs"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

type Command struct {
	Pipelines   []entity.Pipeline
	DryRun      bool
	Healthcheck bool
}

type Handler struct {
	jobs *jobs.Runner
}

type Healthcheck interface {
	Start(ctx context.Context, rid, body string) error
	Success(ctx context.Context, rid, body string) error
	Fail(ctx context.Context, rid, body string) error
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

	hc, err := newHealthcheckService(cmd, pipeline)
	if err != nil {
		return err
	}

	rid := uuid.NewString()

	logHealthcheckErr(ctx, hc.Start(ctx, rid, startBody(pipeline)))

	err = h.jobs.Run(ctx, pipeline.Jobs)
	if err != nil {
		logHealthcheckErr(ctx, hc.Fail(ctx, rid, err.Error()))
		return err
	}

	logHealthcheckErr(ctx, hc.Success(ctx, rid, ""))

	return nil
}

func newHealthcheckService(cmd *Command, pipeline entity.Pipeline) (Healthcheck, error) {
	if !cmd.Healthcheck || cmd.DryRun || strings.TrimSpace(pipeline.HealthcheckURL) == "" {
		return &healthchecks.Dummy{}, nil
	}

	return healthchecks.NewClient(pipeline.HealthcheckURL)
}

func logHealthcheckErr(ctx context.Context, err error) {
	if err == nil {
		return
	}

	log := logger.FromContext(ctx)
	log.Warn().Err(err).Msg("Healthcheck ping failed")
}

func startBody(pipeline entity.Pipeline) string {
	jNames := lo.Map(pipeline.Jobs, func(j entity.Job, _ int) string {
		return j.GetName()
	})

	return fmt.Sprintf("%s: \n\t%s", pipeline.Name, strings.Join(jNames, "\n\t"))
}
