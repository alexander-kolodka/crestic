package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestBackup_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")
	photos := sb.AddRepo("photos")
	cPhotos := sb.AddRepo("cloud_photos")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	sb.AddPipeline("backup_photos").
		Backup("PC", []string{src}, photos).
		Copy("Cloud", photos, cPhotos)

	_, err := sb.Run(ctx, "backup", "--all")

	require.NoError(t, err)
	require.Equal(t, 1, docs.SnapshotCount(ctx))
	require.Equal(t, 1, cDocs.SnapshotCount(ctx))
	require.Equal(t, 1, photos.SnapshotCount(ctx))
	require.Equal(t, 1, cPhotos.SnapshotCount(ctx))
}

func TestBackup_SinglePipeline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")
	photos := sb.AddRepo("photos")
	cPhotos := sb.AddRepo("cloud_photos")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	sb.AddPipeline("backup_photos").
		Backup("PC", []string{src}, photos).
		Copy("Cloud", photos, cPhotos)

	_, err := sb.Run(ctx, "backup", "--pipeline", "backup_docs")

	require.NoError(t, err)
	require.Equal(t, 1, docs.SnapshotCount(ctx))
	require.Equal(t, 1, cDocs.SnapshotCount(ctx))
	require.False(t, photos.IsInitialized(ctx))
	require.False(t, cPhotos.IsInitialized(ctx))
}

func TestBackup_MultiplePipelines(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")
	photos := sb.AddRepo("photos")
	cPhotos := sb.AddRepo("cloud_photos")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	sb.AddPipeline("backup_photos").
		Backup("PC", []string{src}, photos).
		Copy("Cloud", photos, cPhotos)

	_, err := sb.Run(ctx, "backup", "--pipeline", "backup_docs,backup_photos")

	require.NoError(t, err)
	require.Equal(t, 1, docs.SnapshotCount(ctx))
	require.Equal(t, 1, cDocs.SnapshotCount(ctx))
	require.Equal(t, 1, photos.SnapshotCount(ctx))
	require.Equal(t, 1, cPhotos.SnapshotCount(ctx))
}

func TestBackup_SingleJob(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")
	photos := sb.AddRepo("photos")
	cPhotos := sb.AddRepo("cloud_photos")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	sb.AddPipeline("backup_photos").
		Backup("PC", []string{src}, photos).
		Copy("Cloud", photos, cPhotos)

	_, err := sb.Run(ctx, "backup", "--job", "backup_docs/PC")

	require.NoError(t, err)
	require.Equal(t, 1, docs.SnapshotCount(ctx))
	require.False(t, cDocs.IsInitialized(ctx))
	require.False(t, photos.IsInitialized(ctx))
	require.False(t, cPhotos.IsInitialized(ctx))
}

func TestBackup_MultipleJobs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")
	photos := sb.AddRepo("photos")
	cPhotos := sb.AddRepo("cloud_photos")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	sb.AddPipeline("backup_photos").
		Backup("PC", []string{src}, photos).
		Copy("Cloud", photos, cPhotos)

	_, err := sb.Run(ctx, "backup", "--job", "backup_docs/PC,backup_photos/PC")

	require.NoError(t, err)
	require.Equal(t, 1, docs.SnapshotCount(ctx))
	require.Equal(t, 1, photos.SnapshotCount(ctx))
	require.False(t, cDocs.IsInitialized(ctx))
	require.False(t, cPhotos.IsInitialized(ctx))
}

func TestBackup_DryRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	_, err := sb.Run(ctx, "backup", "--all", "--dry-run")

	require.NoError(t, err)
	require.Equal(t, 0, docs.SnapshotCount(ctx))
	require.Equal(t, 0, cDocs.SnapshotCount(ctx))
}

func TestBackup_DryRunSingleJob(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")
	cDocs := sb.AddRepo("cloud_documents")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs).
		Copy("Cloud", docs, cDocs)

	_, err := sb.Run(ctx, "backup", "--job", "backup_docs/PC", "--dry-run")

	require.NoError(t, err)
	require.Equal(t, 0, docs.SnapshotCount(ctx))
	require.False(t, cDocs.IsInitialized(ctx))
}

func TestBackup_FlagConflict(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "all and pipeline",
			args: []string{"--all", "--pipeline", "backup_docs"},
		},
		{
			name: "all and job",
			args: []string{"--all", "--job", "backup_docs/PC"},
		},
		{
			name: "pipeline and job",
			args: []string{"--pipeline", "backup_docs", "--job", "backup_docs/PC"},
		},
		{
			name: "all, pipeline and job",
			args: []string{"--all", "--pipeline", "backup_docs", "--job", "backup_docs/PC"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := harness.New(t)

			_, err := sb.Run(ctx, append([]string{"backup"}, tt.args...)...)

			require.Error(t, err)
			require.Contains(t, err.Error(), "--all, --pipeline and --job are mutually exclusive")
		})
	}
}

func TestBackup_InvalidPipeline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs)

	_, err := sb.Run(ctx, "backup", "--pipeline", "missing")

	require.Error(t, err)
	require.Contains(t, err.Error(), "can't find pipeline missing")
}

func TestBackup_InvalidJob(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)
	src := sb.Mkdir("src")
	sb.WriteFile("src", "file.txt", "hello")

	docs := sb.AddRepo("documents")

	sb.AddPipeline("backup_docs").
		Backup("PC", []string{src}, docs)

	_, err := sb.Run(ctx, "backup", "--job", "backup_docs/unknown")

	require.Error(t, err)
	require.Contains(t, err.Error(), "job unknown not found")
}

func TestBackup_NoFlags(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	sb := harness.New(t)

	_, err := sb.Run(ctx, "backup")

	require.Error(t, err)
	require.Contains(t, err.Error(), "either --all, --pipeline or --job must be specified")
}
