package jobs

import (
	"context"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/hooks"
	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/pkg/mw"
	"github.com/alexander-kolodka/crestic/internal/restic"
)

// Runner executes jobs sequentially (fail-fast).
type Runner struct {
	restic *restic.Service
	hooks  *hooks.Runner
}

// NewRunner creates a job Runner.
func NewRunner(resticService *restic.Service, hooksRunner *hooks.Runner) *Runner {
	return &Runner{
		restic: resticService,
		hooks:  hooksRunner,
	}
}

// Run executes jobs sequentially and returns on the first error (fail-fast).
func (r *Runner) Run(ctx context.Context, jobs []entity.Job) error {
	if isDryRun(ctx) {
		ctx = restic.WithDryRun(ctx)
		ctx = logger.FromContext(ctx).With().Bool("dry-run", true).Logger().WithContext(ctx)
	}

	fn := mw.Chain(
		r.doJob,
		newHookMw(r),
	)

	for _, job := range jobs {
		err := fn(ctx, job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runner) doJob(ctx context.Context, job entity.Job) error {
	switch j := job.(type) {
	case entity.BackupJob:
		jobCtx := logger.WithBackupJobFields(ctx, j)
		log := logger.FromContext(jobCtx)

		err := r.backup(jobCtx, j)
		if err == nil {
			return nil
		}

		log.Error().Msg("Backup job failed")
		return err
	case entity.CopyJob:
		jobCtx := logger.WithCopyJobFields(ctx, j)
		log := logger.FromContext(jobCtx)

		err := r.copy(jobCtx, j)
		if err == nil {
			return nil
		}

		log.Error().Msg("Copy job failed")
		return err
	default:
	}

	return nil
}

func (r *Runner) backup(ctx context.Context, b entity.BackupJob) error {
	ctx = logger.WithBackupJobFields(ctx, b)
	log := logger.FromContext(ctx)
	log.Info().Msg("Processing backup")

	err := r.initRepo(ctx, b.To)
	if err != nil {
		return err
	}

	err = r.restic.Backup(ctx, b)
	if err != nil {
		return err
	}

	err = r.restic.Check(ctx, b.To)
	if err != nil {
		return err
	}

	err = r.restic.Forget(ctx, b.To)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) copy(ctx context.Context, c entity.CopyJob) error {
	log := logger.FromContext(ctx)
	log.Info().Msg("Processing copy")

	err := r.initRepo(ctx, c.From)
	if err != nil {
		return err
	}

	err = r.initRepo(ctx, c.To)
	if err != nil {
		return err
	}

	err = r.restic.Copy(ctx, c)
	if err != nil {
		return err
	}

	err = r.restic.Check(ctx, c.To)
	if err != nil {
		return err
	}

	err = r.restic.Forget(ctx, c.To)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runner) initRepo(ctx context.Context, repo *entity.Repository) error {
	isRepoInitialized, err := r.restic.IsRepoInitialized(ctx, repo)
	if err != nil {
		return err
	}

	if isRepoInitialized {
		return nil
	}

	return r.restic.Init(ctx, repo)
}

func (r *Runner) executeHooks(ctx context.Context, phase string, cmds []string) error {
	return r.hooks.Execute(ctx, phase, cmds)
}
