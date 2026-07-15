package jobs

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/hooks"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/mw"
)

type hookExecutor interface {
	executeHooks(ctx context.Context, phase string, hooks []string) error
}

func newHookMw(h hookExecutor) mw.Middleware[entity.Job] {
	return func(fn mw.Func[entity.Job]) mw.Func[entity.Job] {
		return func(ctx context.Context, j entity.Job) error {
			jobHooks := j.GetHooks()
			jName := j.GetFullName()

			err := h.executeHooks(hooks.WithJobEnv(ctx, jName, nil), hooks.PhaseBefore, jobHooks.Before)
			if err != nil {
				logFailureHookErr(
					ctx,
					h.executeHooks(hooks.WithJobEnv(ctx, jName, err), hooks.PhaseFailure, jobHooks.Failure),
				)
				return fmt.Errorf("before hooks failed: %w", err)
			}

			err = fn(ctx, j)
			if err != nil {
				logFailureHookErr(
					ctx,
					h.executeHooks(hooks.WithJobEnv(ctx, jName, err), hooks.PhaseFailure, jobHooks.Failure),
				)
				return err
			}

			return h.executeHooks(hooks.WithJobEnv(ctx, jName, nil), hooks.PhaseSuccess, jobHooks.Success)
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
