package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestCron_RunsDuePipeline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").Cron("* * * * *").Backup("local", []string{src}, repo)
	sb.WriteCronState("nightly", time.Now().Add(-2*time.Hour))

	require.False(t, repo.IsInitialized(ctx))

	_, err := sb.Run(ctx, "cron")
	require.NoError(t, err)
	require.Equal(t, 1, repo.SnapshotCount(ctx))
}

func TestCron_RunsPipelineOnFirstRunWithinGrace(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").Cron("* * * * *").Backup("local", []string{src}, repo)

	require.False(t, repo.IsInitialized(ctx))

	_, err := sb.Run(ctx, "cron")
	require.NoError(t, err)
	require.Equal(t, 1, repo.SnapshotCount(ctx))
}
