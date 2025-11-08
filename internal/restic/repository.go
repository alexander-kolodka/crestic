package restic

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

const (
	// Exit status is 3 if some source data could not be read (incomplete snapshot created)
	// This happens with Cryptomator or other FUSE mounts, where xattrs are inaccessible.
	// Despite this warning, the backup is completed successfully, so this exit code could be ignored.
	resticBackupExitCodeMissingXAttrs = 3

	// Exit code is 10 if the repository does not exist.
	resticStatsExitCodeRepoDoesNotExist = 10
)

// Service provides high-level operations for interacting with restic repositories.
type Service struct {
	runner runner
}

func NewService(runner runner) *Service {
	return &Service{
		runner: &resticRunner{runner: runner},
	}
}

// Init initializes a new restic repository at the specified path.
// This must be called before the repository can be used for backups.
func (r *Service) Init(ctx context.Context, repo *entity.Repository) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Initializing repository")

	result := r.runner.Run(
		ctx,
		"restic",
		"init",
		"-r", repo.Path,
		"--password-command", repo.PasswordCMD,
	)

	return r.toErr(ctx, result, repo, "init")
}

// IsRepoInitialized checks if a repository has been initialized and is accessible.
func (r *Service) IsRepoInitialized(ctx context.Context, repo *entity.Repository) (bool, error) {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Checking if repository is initialized")

	result := r.runner.Run(
		shell.WithSilence(ctx),
		"restic",
		"stats",
		"-r", repo.Path,
		"--password-command", repo.PasswordCMD,
	)

	if result.ExitCode == resticStatsExitCodeRepoDoesNotExist {
		return false, nil
	}

	if result.Error == nil {
		return true, nil
	}

	return false, r.toErr(ctx, result, repo, "stats")
}

// Backup creates a new backup snapshot from the specified source directories.
// If the context contains a dry-run flag, no actual backup is performed.
func (r *Service) Backup(ctx context.Context, b entity.BackupJob) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Starting backup")

	args := []string{
		"backup",
		"-r", b.To.Path,
		"--password-command", b.To.PasswordCMD,
	}

	if isDryRun(ctx) {
		args = append(args, "--dry-run")
	}

	args = append(args, b.Options.ToArgs()...)
	args = append(args, b.From...)

	result := r.runner.Run(ctx, "restic", args...)

	if b.IgnoreMissingXAttrsError && result.ExitCode == resticBackupExitCodeMissingXAttrs {
		log.Warn().Msg("Backup failed with missing xattrs, but it was ignored")
		return nil
	}

	return r.toErr(ctx, result, b.To, "backup")
}

// Check verifies the integrity of a repository.
func (r *Service) Check(ctx context.Context, repo *entity.Repository) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Running integrity check")

	result := r.runner.Run(
		ctx,
		"restic",
		"check",
		"-r", repo.Path,
		"--password-command", repo.PasswordCMD,
	)

	return r.toErr(ctx, result, repo, "check")
}

// Forget removes old snapshots according to the repository's retention policy.
// This only marks snapshots for deletion; use with --prune flag (in ForgetOptions)
// to actually remove the data and free disk space.
func (r *Service) Forget(ctx context.Context, repo *entity.Repository) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Running forget")

	args := []string{
		"forget",
		"-r", repo.Path,
		"--password-command", repo.PasswordCMD,
	}

	if isDryRun(ctx) {
		args = append(args, "--dry-run")
	}

	args = append(args, repo.ForgetOptions.ToArgs()...)

	result := r.runner.Run(ctx, "restic", args...)

	return r.toErr(ctx, result, repo, "forget")
}

// Copy copies snapshots from one repository to another.
func (r *Service) Copy(ctx context.Context, job entity.CopyJob) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Starting copy")

	args := []string{
		"copy",
		"-r", job.To.Path,
		"--from-repo", job.From.Path,
		"--password-command", job.To.PasswordCMD,
		"--from-password-command", job.From.PasswordCMD,
	}

	args = append(args, job.Options.ToArgs()...)

	if isDryRun(ctx) {
		log.Debug().
			Strs("args", args).
			Msg("DRY RUN: would execute restic copy")
		return nil
	}

	result := r.runner.Run(ctx, "restic", args...)

	if result.Error == nil {
		return nil
	}

	return fmt.Errorf(
		"repository %s: restic copy from %s failed [exit code %d]: %w",
		job.To.Name,
		job.From.Name,
		result.ExitCode,
		result.Error,
	)
}

// Restore extracts files from a snapshot to the specified target directory.
// The snapshot parameter can be a snapshot ID or "latest" for the most recent snapshot.
func (r *Service) Restore(
	ctx context.Context,
	repo *entity.Repository,
	target, snapshot string,
) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Starting restore")

	args := []string{
		"restore",
		"--target", target,
		"-r", repo.Path,
		"--password-command", repo.PasswordCMD,
		snapshot,
	}

	result := r.runner.Run(ctx, "restic", args...)
	return r.toErr(ctx, result, repo, "restore")
}

// Exec executes an arbitrary restic command on a repository.
func (r *Service) Exec(
	ctx context.Context,
	repo *entity.Repository,
	cmd string,
	args []string,
) error {
	log := logger.FromContext(ctx)
	log.Debug().Msg("Executing restic command")

	args = append(
		[]string{
			cmd,
			"-r", repo.Path,
			"--password-command", repo.PasswordCMD,
		},
		args...,
	)

	result := r.runner.Run(ctx, "restic", args...)

	return r.toErr(ctx, result, repo, cmd)
}

// Unlock removes stale locks from a repository.
func (r *Service) Unlock(ctx context.Context, repo *entity.Repository) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Unlocking repository")

	result := r.runner.Run(
		ctx,
		"restic",
		"unlock",
		"-r", repo.Path,
		"--password-command", repo.PasswordCMD,
	)

	return r.toErr(ctx, result, repo, "unlock")
}

func (r *Service) toErr(ctx context.Context, result *shell.Result, repo *entity.Repository, cmdName string) error {
	if result.Error == nil {
		return nil
	}

	log := logger.FromContext(ctx)
	log.Error().
		Str("cmd", cmdName).
		Int("exit_code", result.ExitCode).
		Err(result.Error).
		Msg("restic command failed")

	return fmt.Errorf(
		"repository %s: restic %s failed [exit code %d]: %w",
		repo.Name,
		cmdName,
		result.ExitCode,
		result.Error,
	)
}
