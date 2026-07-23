package mw_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/internal/pkg/mw"
	"github.com/alexander-kolodka/crestic/internal/pkg/testutils"
)

func TestChain(t *testing.T) {
	var stack []string

	base := func(_ context.Context, _ string) error {
		stack = append(stack, "base")
		return nil
	}

	mw1 := func(next mw.Func[string]) mw.Func[string] {
		return func(ctx context.Context, s string) error {
			stack = append(stack, "mw1: before")
			err := next(ctx, s)
			stack = append(stack, "mw1: after")
			return err
		}
	}

	mw2 := func(next mw.Func[string]) mw.Func[string] {
		return func(ctx context.Context, s string) error {
			stack = append(stack, "mw2: before")
			err := next(ctx, s)
			stack = append(stack, "mw2: after")
			return err
		}
	}

	chained := mw.Chain(base, mw1, mw2)

	require.NoError(t, chained(context.Background(), "x"))

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
