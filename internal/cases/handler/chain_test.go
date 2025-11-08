package handler_test

import (
	"context"
	"testing"

	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	type testCommand struct{}

	var stack []string

	base := handler.NewHandler(func(_ context.Context, _ testCommand) error {
		stack = append(stack, "base")
		return nil
	})

	mw1 := func(h handler.Handler[testCommand]) handler.Handler[testCommand] {
		return handler.NewHandler(func(ctx context.Context, cmd testCommand) error {
			stack = append(stack, "mw1: before")
			err := h.Handle(ctx, cmd)
			stack = append(stack, "mw1: after")
			return err
		})
	}

	mw2 := func(h handler.Handler[testCommand]) handler.Handler[testCommand] {
		return handler.NewHandler(func(ctx context.Context, cmd testCommand) error {
			stack = append(stack, "mw2: before")
			err := h.Handle(ctx, cmd)
			stack = append(stack, "mw2: after")
			return err
		})
	}

	chained := handler.Chain(base, mw1, mw2)

	require.NoError(t, chained.Handle(context.Background(), testCommand{}))

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
