package backup

import "C"

import (
	"context"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

type Command struct {
	Jobs   []entity.Job
	DryRun bool
}

type Handler struct {
	restic *restic.Service
	runner *shell.Executor
	hc     *healthchecks.Client
}

// NewHandler creates a backup command Handler.
func NewHandler(restic *restic.Service, runner *shell.Executor, hc *healthchecks.Client) *Handler {
	return &Handler{
		restic: restic,
		runner: runner,
		hc:     hc,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd *Command) error {
	if cmd.DryRun {
		ctx = restic.WithDryRun(ctx)
		ctx = logger.FromContext(ctx).With().Bool("dry-run", cmd.DryRun).Logger().WithContext(ctx)
	}

	fn := chain(
		h.doJob,
		newHealthcheckMw(h.hc),
		newHookMw(h),
	)

	jobErrors := newJobErrors()

	for _, job := range cmd.Jobs {
		err := fn(ctx, job)
		if err != nil {
			jobErrors.Add(job.GetName(), err)
			continue
		}
	}

	if jobErrors.HasErrors() {
		return jobErrors
	}

	return nil
}

func (h *Handler) doJob(ctx context.Context, job entity.Job) error {
	switch j := job.(type) {
	case entity.BackupJob:
		jobCtx := logger.WithBackupJobFields(ctx, j)
		log := logger.FromContext(jobCtx)

		err := h.backup(jobCtx, j)
		if err == nil {
			return nil
		}

		log.Error().Msg("Backup job failed")
		return err
	case entity.CopyJob:
		jobCtx := logger.WithCopyJobFields(ctx, j)
		log := logger.FromContext(jobCtx)

		err := h.copy(jobCtx, j)
		if err == nil {
			return nil
		}

		log.Error().Msg("Copy job failed")
		return err
	default:
	}

	return nil
}

func (h *Handler) backup(ctx context.Context, b entity.BackupJob) error {
	ctx = logger.WithBackupJobFields(ctx, b)
	log := logger.FromContext(ctx)
	log.Info().Msg("Processing backup")

	err := h.initRepo(ctx, b.To)
	if err != nil {
		return err
	}

	err = h.restic.Backup(ctx, b)
	if err != nil {
		return err
	}

	err = h.restic.Check(ctx, b.To)
	if err != nil {
		return err
	}

	err = h.restic.Forget(ctx, b.To)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) copy(ctx context.Context, c entity.CopyJob) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Processing copy")

	err := h.initRepo(ctx, c.From)
	if err != nil {
		return err
	}

	err = h.initRepo(ctx, c.To)
	if err != nil {
		return err
	}

	err = h.restic.Copy(ctx, c)
	if err != nil {
		return err
	}

	err = h.restic.Check(ctx, c.To)
	if err != nil {
		return err
	}

	err = h.restic.Forget(ctx, c.To)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) initRepo(ctx context.Context, repo *entity.Repository) error {
	isRepoInitialized, err := h.restic.IsRepoInitialized(ctx, repo)
	if err != nil {
		return err
	}

	if isRepoInitialized {
		return nil
	}

	return h.restic.Init(ctx, repo)
}

func (h *Handler) executeHooks(ctx context.Context, hooks []string) error {
	ctx = logger.WithSource(ctx, "hooks")
	for _, hook := range hooks {
		result := h.runner.Run(ctx, "sh", "-c", hook)
		if result.Error != nil {
			return fmt.Errorf(
				`hook failed "%s" [exit code %d]: %w`,
				hook,
				result.ExitCode,
				result.Error,
			)
		}
	}
	return nil
}
