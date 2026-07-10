package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestForget(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	repo := sb.AddRepo("documents")
	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, repo)

	_, err := sb.Run(ctx, "backup", "--all")
	require.NoError(t, err)

	_, err = sb.Run(ctx, "forget", "--repo", repo.Name)
	require.NoError(t, err)
}
