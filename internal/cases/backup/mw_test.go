package backup

import (
	"context"
	"testing"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	var stack []string

	base := func(_ context.Context, _ entity.Job) error {
		stack = append(stack, "base")
		return nil
	}

	mw1 := func(base do) do {
		return func(ctx context.Context, j entity.Job) error {
			stack = append(stack, "mw1: before")
			err := base(ctx, j)
			stack = append(stack, "mw1: after")
			return err
		}
	}

	mw2 := func(base do) do {
		return func(ctx context.Context, j entity.Job) error {
			stack = append(stack, "mw2: before")
			err := base(ctx, j)
			stack = append(stack, "mw2: after")
			return err
		}
	}

	chained := chain(base, mw1, mw2)

	require.NoError(t, chained(context.Background(), entity.BackupJob{}))

	expected := []string{
		"mw1: before",
		"mw2: before",
		"base",
		"mw2: after",
		"mw1: after",
	}
	testutils.Equal(
		t,
		expected,
		stack,
	)
}
