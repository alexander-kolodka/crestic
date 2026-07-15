package runpipelines

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/mw"
)

func newHealthcheckMw(cmd *Command) mw.Middleware[entity.Pipeline] {
	return func(fn mw.Func[entity.Pipeline]) mw.Func[entity.Pipeline] {
		return func(ctx context.Context, pipeline entity.Pipeline) error {
			hc, err := newHealthcheckService(cmd, pipeline)
			if err != nil {
				return err
			}

			rid := uuid.NewString()
			logHealthcheckErr(ctx, hc.Start(ctx, rid, startBody(pipeline)))

			err = fn(ctx, pipeline)
			if err != nil {
				logHealthcheckErr(ctx, hc.Fail(ctx, rid, err.Error()))
				return err
			}

			logHealthcheckErr(ctx, hc.Success(ctx, rid, ""))
			return nil
		}
	}
}

//nolint:ireturn // factory returns Dummy or Client behind the local Healthcheck interface
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
