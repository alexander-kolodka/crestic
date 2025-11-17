package logger

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/alexander-kolodka/crestic/internal/entity"
)

// FromContext extracts logger from context or returns global logger.
func FromContext(ctx context.Context) zerolog.Logger {
	return *zerolog.Ctx(ctx)
}

// WithBackupJobFields adds job-related fields to context for backup job.
func WithBackupJobFields(ctx context.Context, job entity.BackupJob) context.Context {
	return FromContext(ctx).With().
		Str("job", job.Name).
		Str("repo", job.To.Name).
		Str("repo_path", job.To.Path).
		Strs("backup_sources", job.From).
		Logger().WithContext(ctx)
}

// WithCopyJobFields adds job-related fields to context for copy job.
func WithCopyJobFields(ctx context.Context, job entity.CopyJob) context.Context {
	return FromContext(ctx).With().
		Str("job", job.Name).
		Str("from_repo", job.From.Name).
		Str("from_repo_path", job.From.Path).
		Str("to_repo", job.To.Name).
		Str("to_repo_path", job.To.Path).
		Logger().WithContext(ctx)
}

// WithRepoFields adds repository fields to context.
func WithRepoFields(ctx context.Context, repo *entity.Repository) context.Context {
	return FromContext(ctx).With().
		Str("repo", repo.Name).
		Str("repo_path", repo.Path).
		Logger().WithContext(ctx)
}

type jsonModeKey struct{}

// WithJSONMode adds JSON mode flag to context.
func WithJSONMode(ctx context.Context) context.Context {
	return context.WithValue(ctx, jsonModeKey{}, true)
}

// IsJSONMode checks if JSON mode is enabled in context.
func IsJSONMode(ctx context.Context) bool {
	json, ok := ctx.Value(jsonModeKey{}).(bool)
	return ok && json
}

type sourceKey struct{}

// WithSource adds source field to the context logger.
func WithSource(ctx context.Context, source string) context.Context {
	ctx = context.WithValue(ctx, sourceKey{}, source)
	return FromContext(ctx).With().
		Str("source", source).
		Logger().WithContext(ctx)
}

// GetSource extracts source field from context.
func GetSource(ctx context.Context) string {
	source, ok := ctx.Value(sourceKey{}).(string)
	if ok {
		return source
	}
	return ""
}
