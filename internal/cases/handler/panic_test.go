package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/panix"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithPanicRecoveryCatchesPanic(t *testing.T) {
	type testCommand struct{}

	base := handler.NewHandler(func(_ context.Context, _ testCommand) error {
		panic("test panic")
	})

	wrapped := handler.WithPanicRecovery[testCommand]()(base)

	err := wrapped.Handle(context.Background(), testCommand{})
	var pe *panix.PanicError
	assert.ErrorAs(t, err, &pe)
}

func TestWithPanicRecoveryWithError(t *testing.T) {
	type testCommand struct{}

	testErr := errors.New("test error")
	base := handler.NewHandler(func(_ context.Context, _ testCommand) error {
		return testErr
	})

	wrapped := handler.WithPanicRecovery[testCommand]()(base)
	err := wrapped.Handle(context.Background(), testCommand{})
	require.ErrorIs(t, err, testErr)
	var pe *panix.PanicError
	assert.NotErrorAs(t, err, &pe)
}

func TestWithPanicRecoveryWithoutPanic(t *testing.T) {
	type testCommand struct{}

	base := handler.NewHandler(func(_ context.Context, _ testCommand) error {
		return nil
	})

	wrapped := handler.WithPanicRecovery[testCommand]()(base)

	err := wrapped.Handle(context.Background(), testCommand{})
	assert.NoError(t, err)
}
