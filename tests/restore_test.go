package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestRestore_Latest(t *testing.T) {
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

	target := sb.Mkdir("restore")
	_, err = sb.Run(ctx, "restore", "--repo", repo.Name, "--target", target)
	require.NoError(t, err)

	restored := harness.FindFile(t, target, "file.txt")
	content, err := os.ReadFile(restored)

	require.NoError(t, err)
	require.Equal(t, "hello", string(content))
}
