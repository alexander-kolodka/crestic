package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestUnlock_Repo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("documents")

	_, err := sb.Run(ctx, "check", "--repo", repo.Name)
	require.NoError(t, err)

	_, err = sb.Run(ctx, "unlock", "--repo", repo.Name)
	require.NoError(t, err)
}
