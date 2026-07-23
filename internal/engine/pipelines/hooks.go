package pipelines

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/engine/hooks"
	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/pkg/mw"
)

func newHookMw(runner *hooks.Runner) mw.Middleware[entity.Pipeline] {
	return func(fn mw.Func[entity.Pipeline]) mw.Func[entity.Pipeline] {
		return func(ctx context.Context, pipeline entity.Pipeline) error {
			pHooks := pipeline.Hooks
			pName := pipeline.Name

			err := runner.Execute(hooks.WithPipelineEnv(ctx, pName, nil), hooks.PhaseBefore, pHooks.Before)
			if err != nil {
				logFailureHookErr(
					ctx,
					runner.Execute(hooks.WithPipelineEnv(ctx, pName, err), hooks.PhaseFailure, pHooks.Failure),
				)
				return fmt.Errorf("before hooks failed: %w", err)
			}

			err = fn(ctx, pipeline)
			if err != nil {
				logFailureHookErr(
					ctx,
					runner.Execute(hooks.WithPipelineEnv(ctx, pName, err), hooks.PhaseFailure, pHooks.Failure),
				)
				return err
			}

			return runner.Execute(hooks.WithPipelineEnv(ctx, pName, nil), hooks.PhaseSuccess, pHooks.Success)
		}
	}
}

func logFailureHookErr(ctx context.Context, err error) {
	if err == nil {
		return
	}

	log := logger.FromContext(ctx)
	log.Warn().Err(err).Msg("Failure hook failed")
}
