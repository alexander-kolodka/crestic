package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestCheck_Repo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("documents")

	require.False(t, repo.IsInitialized(ctx))

	_, err := sb.Run(ctx, "check", "--repo", repo.Name)
	require.NoError(t, err)
	require.True(t, repo.IsInitialized(ctx))
}

func TestCheck_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	docs := sb.AddRepo("documents")
	photos := sb.AddRepo("photos")

	require.False(t, docs.IsInitialized(ctx))
	require.False(t, photos.IsInitialized(ctx))

	_, err := sb.Run(ctx, "check", "--all")
	require.NoError(t, err)
	require.True(t, docs.IsInitialized(ctx))
	require.True(t, photos.IsInitialized(ctx))
}
