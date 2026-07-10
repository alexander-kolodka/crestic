package harness

import (
	"context"
	"encoding/json"

	"github.com/stretchr/testify/require"
)

// RepoDir is a restic repository directory registered in the sandbox config.
type RepoDir struct {
	sandbox *Sandbox
	Name    string
	Path    string
}

// IsInitialized reports whether the repository has been initialized with restic.
func (r *RepoDir) IsInitialized(ctx context.Context) bool {
	r.sandbox.t.Helper()

	_, err := r.sandbox.Run(ctx, "exec", "--repo", r.Name, "stats")
	return err == nil
}

// SnapshotCount returns the number of snapshots in the repository.
func (r *RepoDir) SnapshotCount(ctx context.Context) int {
	r.sandbox.t.Helper()

	output, err := r.sandbox.Run(ctx, "exec", "--repo", r.Name, "snapshots", "--json")
	require.NoError(r.sandbox.t, err, output)

	var snapshots []json.RawMessage
	require.NoError(r.sandbox.t, json.Unmarshal([]byte(output), &snapshots))
	return len(snapshots)
}
