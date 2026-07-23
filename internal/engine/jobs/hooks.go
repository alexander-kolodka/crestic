package jobs

import (
	"context"
	"fmt"

	hooks2 "github.com/alexander-kolodka/crestic/internal/engine/hooks"
	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/pkg/mw"
)

type hookExecutor interface {
	executeHooks(ctx context.Context, phase string, hooks []string) error
}

func newHookMw(h hookExecutor) mw.Middleware[entity.Job] {
	return func(fn mw.Func[entity.Job]) mw.Func[entity.Job] {
		return func(ctx context.Context, j entity.Job) error {
			jobHooks := j.GetHooks()
			jName := j.GetFullName()

			err := h.executeHooks(hooks2.WithJobEnv(ctx, jName, nil), hooks2.PhaseBefore, jobHooks.Before)
			if err != nil {
				logFailureHookErr(
					ctx,
					h.executeHooks(hooks2.WithJobEnv(ctx, jName, err), hooks2.PhaseFailure, jobHooks.Failure),
				)
				return fmt.Errorf("before hooks failed: %w", err)
			}

			err = fn(ctx, j)
			if err != nil {
				logFailureHookErr(
					ctx,
					h.executeHooks(hooks2.WithJobEnv(ctx, jName, err), hooks2.PhaseFailure, jobHooks.Failure),
				)
				return err
			}

			return h.executeHooks(hooks2.WithJobEnv(ctx, jName, nil), hooks2.PhaseSuccess, jobHooks.Success)
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
