package backup

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
)

// newHealthcheckMw wraps a do func with Healthchecks.io monitoring.
// Automatically sends start/success/fail signals with unique run ID for grouping.
// Signal sending errors don't abort the main task execution.
func newHealthcheckMw(hc *healthchecks.Client) mw {
	return func(fn do) do {
		return func(ctx context.Context, j entity.Job) error {
			url := j.GetHealthcheckURL()
			if url == "" {
				return fn(ctx, j)
			}

			rid := uuid.NewString()
			p := healthchecks.Payload{JobName: j.GetName()}

			_ = hc.Start(ctx, url, rid, p)

			err := fn(ctx, j)
			if err != nil {
				p.Err = err.Error()
				_ = hc.Fail(ctx, url, rid, p)
				return err
			}

			_ = hc.Success(ctx, url, rid, p)
			return nil
		}
	}
}
