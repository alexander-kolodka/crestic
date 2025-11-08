package restic

import (
	"context"
)

type dryRun struct{}

func WithDryRun(ctx context.Context) context.Context {
	return context.WithValue(ctx, dryRun{}, true)
}

func isDryRun(ctx context.Context) bool {
	dry, ok := ctx.Value(dryRun{}).(bool)
	return ok && dry
}
