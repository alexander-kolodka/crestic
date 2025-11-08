package backup

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

type hookExecutor interface {
	executeHooks(ctx context.Context, hooks []string) error
}

func newHookMw(h hookExecutor) mw {
	return func(fn do) do {
		return func(ctx context.Context, j entity.Job) error {
			hooks := j.GetHooks()
			jName := j.GetName()

			err := h.executeHooks(withEnv(ctx, jName, nil), hooks.Before)
			if err != nil {
				_ = h.executeHooks(withEnv(ctx, jName, err), hooks.Failure)
				return fmt.Errorf("before hooks failed: %w", err)
			}

			err = fn(ctx, j)
			if err != nil {
				_ = h.executeHooks(withEnv(ctx, jName, err), hooks.Failure)
				return err
			}

			return h.executeHooks(withEnv(ctx, jName, nil), hooks.Success)
		}
	}
}

func withEnv(ctx context.Context, jobName string, err error) context.Context {
	env := map[string]string{
		"CRESTIC_JOB_NAME": jobName,
	}

	if err != nil {
		env["CRESTIC_ERROR"] = err.Error()
	}

	return shell.WithEnv(ctx, env)
}
