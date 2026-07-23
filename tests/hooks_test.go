package tests

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/internal/pkg/testutils"
	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestHooks_BackupHappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").
		Hooks(harness.Hooks{
			Before:  []string{appendHookOrder("p.before")},
			Success: []string{appendHookOrder("p.success")},
			Failure: []string{appendHookOrder("p.failure")},
		}).
		Backup("local", []string{src}, repo).
		JobHooks(harness.Hooks{
			Before:  []string{appendHookOrder("j.before")},
			Success: []string{appendHookOrder("j.success")},
			Failure: []string{appendHookOrder("j.failure")},
		})

	_, err := sb.Run(ctx, "backup", "--all")
	require.NoError(t, err)
	require.Equal(t, 1, repo.SnapshotCount(ctx))
	requireHookOrder(t, sb, "p.before", "j.before", "j.success", "p.success")
}

func TestHooks_JobBeforeFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").
		Hooks(harness.Hooks{
			Before:  []string{appendHookOrder("p.before")},
			Success: []string{appendHookOrder("p.success")},
			Failure: []string{appendHookOrder("p.failure")},
		}).
		Backup("local", []string{src}, repo).
		JobHooks(harness.Hooks{
			Before:  []string{appendHookOrder("j.before"), "false"},
			Success: []string{appendHookOrder("j.success")},
			Failure: []string{appendHookOrder("j.failure")},
		})

	_, err := sb.Run(ctx, "backup", "--all")
	require.Error(t, err)
	require.False(t, repo.IsInitialized(ctx))
	requireHookOrder(t, sb, "p.before", "j.before", "j.failure", "p.failure")
}

func TestHooks_JobBodyFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	missing := sb.Path("missing-src")

	sb.AddPipeline("nightly").
		Hooks(harness.Hooks{
			Before:  []string{appendHookOrder("p.before")},
			Success: []string{appendHookOrder("p.success")},
			Failure: []string{appendHookOrder("p.failure")},
		}).
		Backup("local", []string{missing}, repo).
		JobHooks(harness.Hooks{
			Before:  []string{appendHookOrder("j.before")},
			Success: []string{appendHookOrder("j.success")},
			Failure: []string{appendHookOrder("j.failure")},
		})

	_, err := sb.Run(ctx, "backup", "--all")
	require.Error(t, err)
	require.Equal(t, 0, repo.SnapshotCount(ctx))
	requireHookOrder(t, sb, "p.before", "j.before", "j.failure", "p.failure")
}

func TestHooks_PipelineBeforeFails(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").
		Hooks(harness.Hooks{
			Before:  []string{appendHookOrder("p.before"), "false"},
			Success: []string{appendHookOrder("p.success")},
			Failure: []string{appendHookOrder("p.failure")},
		}).
		Backup("local", []string{src}, repo).
		JobHooks(harness.Hooks{
			Before:  []string{appendHookOrder("j.before")},
			Success: []string{appendHookOrder("j.success")},
			Failure: []string{appendHookOrder("j.failure")},
		})

	_, err := sb.Run(ctx, "backup", "--all")
	require.Error(t, err)
	require.False(t, repo.IsInitialized(ctx))
	requireHookOrder(t, sb, "p.before", "p.failure")
}

func TestHooks_BackupJobSkipsPipelineHooks(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").
		Hooks(harness.Hooks{
			Before:  []string{appendHookOrder("p.before")},
			Success: []string{appendHookOrder("p.success")},
			Failure: []string{appendHookOrder("p.failure")},
		}).
		Backup("local", []string{src}, repo).
		JobHooks(harness.Hooks{
			Before:  []string{appendHookOrder("j.before")},
			Success: []string{appendHookOrder("j.success")},
			Failure: []string{appendHookOrder("j.failure")},
		})

	_, err := sb.Run(ctx, "backup", "--job", "nightly/local")
	require.NoError(t, err)
	require.Equal(t, 1, repo.SnapshotCount(ctx))
	requireHookOrder(t, sb, "j.before", "j.success")
}

func TestHooks_CronHappyPath(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	repo := sb.AddRepo("docs")
	src := sb.Mkdir("data")
	sb.WriteFile("data", "file.txt", "hello")

	sb.AddPipeline("nightly").
		Cron("* * * * *").
		Hooks(harness.Hooks{
			Before:  []string{appendHookOrder("p.before")},
			Success: []string{appendHookOrder("p.success")},
			Failure: []string{appendHookOrder("p.failure")},
		}).
		Backup("local", []string{src}, repo).
		JobHooks(harness.Hooks{
			Before:  []string{appendHookOrder("j.before")},
			Success: []string{appendHookOrder("j.success")},
			Failure: []string{appendHookOrder("j.failure")},
		})
	sb.WriteCronState("nightly", time.Now().Add(-2*time.Hour))

	_, err := sb.Run(ctx, "cron")
	require.NoError(t, err)
	require.Equal(t, 1, repo.SnapshotCount(ctx))
	requireHookOrder(t, sb, "p.before", "j.before", "j.success", "p.success")
}

func appendHookOrder(step string) string {
	return "echo " + step + " >> hook-order"
}

func requireHookOrder(t *testing.T, sb *harness.Sandbox, want ...string) {
	t.Helper()

	data, err := os.ReadFile(sb.Path("hook-order"))
	require.NoError(t, err)

	got := strings.Split(strings.TrimSpace(string(data)), "\n")
	testutils.Equal(t, want, got)
}
