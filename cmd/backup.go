package cmd

import (
	"errors"
	"fmt"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/alexander-kolodka/crestic/internal/cases/backup"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/cases/runpipelines"
	"github.com/alexander-kolodka/crestic/internal/engine/hooks"
	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/engine/restic"
	"github.com/alexander-kolodka/crestic/internal/engine/shell"
	"github.com/alexander-kolodka/crestic/internal/entity"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create encrypted backups of configured sources",
	Long: `Create encrypted backups of directories specified in your configuration file.

This command performs a complete backup workflow for each job:

The backup process (automatic):
  1. Runs 'before' hooks (if configured)
  3. Checks if repository is initialized (auto-initializes if needed)
  4. Creates encrypted backup snapshot using restic
  5. Verifies repository integrity (restic check)
  6. Applies retention policy (restic forget with forget_options)
  7. Runs 'success' or 'failure' hooks based on outcome

Note: The forget step automatically runs after each backup if forget_options
are configured in the repository. If --prune flag is set in forget_options,
old data is actually removed from the repository to free disk space.

Jobs run sequentially. If a job fails, the remaining jobs in the list are not
executed and the error is returned immediately.

Examples:
  # Run all jobs from all pipelines
  crestic backup --all

  # Run all jobs in a pipeline
  crestic backup --pipeline documents

  # Run a specific job
  crestic backup --job documents/local-backup

  # Dry run (show what would be backed up)
  crestic backup --all --dry-run`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		err = validatePipelinesAndJobs(cmd)
		if err != nil {
			return err
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		healthcheck, _ := cmd.Flags().GetBool("healthcheck")

		executor := shell.NewExecutor()
		hooksRunner := hooks.New(executor)
		jobRunner := jobs.NewRunner(restic.NewService(executor), hooksRunner)

		jobsHandler := handler.Chain(
			backup.NewHandler(jobRunner),
			handler.WithPanicRecovery[*backup.Command](),
		)

		fullJobNames, _ := cmd.Flags().GetStringSlice("job")
		extractedJobs, err := extractJobs(cfg, fullJobNames)
		if err != nil {
			return err
		}

		if len(extractedJobs) > 0 {
			return jobsHandler.Handle(cmd.Context(), &backup.Command{
				Jobs:   extractedJobs,
				DryRun: dryRun,
			})
		}

		runPipelinesHandler := handler.Chain(
			runpipelines.NewHandler(jobRunner, hooksRunner),
			handler.WithPanicRecovery[*runpipelines.Command](),
		)

		all, _ := cmd.Flags().GetBool("all")
		if all {
			return runPipelinesHandler.Handle(cmd.Context(), &runpipelines.Command{
				Pipelines:   cfg.Pipelines,
				DryRun:      dryRun,
				Healthcheck: healthcheck,
			})
		}

		pipelineNames, _ := cmd.Flags().GetStringSlice("pipeline")
		pipelines, err := extractPipelines(cfg, pipelineNames)
		if err != nil {
			return err
		}

		return runPipelinesHandler.Handle(cmd.Context(), &runpipelines.Command{
			Pipelines:   pipelines,
			DryRun:      dryRun,
			Healthcheck: healthcheck,
		})
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().BoolP("all", "a", false, "Run all jobs from all pipelines")
	backupCmd.Flags().
		StringSliceP("job", "j", nil, "Run specific jobs by qualified name pipeline/job (comma-separated)")
	backupCmd.Flags().StringSliceP("pipeline", "p", nil, "Run all jobs in specific pipelines by name (comma-separated)")
	backupCmd.Flags().Bool("dry-run", false, "Dry run")
	backupCmd.Flags().Bool("healthcheck", false, "Send healthcheck pings for pipelines with healthcheck_url")

	_ = backupCmd.RegisterFlagCompletionFunc("job", jobAutocompletion)
	_ = backupCmd.RegisterFlagCompletionFunc("pipeline", pipelineAutocompletion)
}

func extractJobs(cfg *entity.Config, fullJobNames []string) ([]entity.Job, error) {
	return lo.MapErr(fullJobNames, func(name string, _ int) (entity.Job, error) {
		return cfg.FindJob(name)
	})
}

func extractPipelines(cfg *entity.Config, pipelines []string) ([]entity.Pipeline, error) {
	return lo.MapErr(pipelines, func(name string, _ int) (entity.Pipeline, error) {
		p, ok := cfg.FindPipeline(name)
		if !ok {
			return entity.Pipeline{}, fmt.Errorf("can't find pipeline %s", name)
		}

		return p, nil
	})
}

func validatePipelinesAndJobs(cmd *cobra.Command) error {
	err := assertNoPipelinesJobsConflict(cmd)
	if err != nil {
		return err
	}

	return nil
}

func assertNoPipelinesJobsConflict(cmd *cobra.Command) error {
	all, _ := cmd.Flags().GetBool("all")
	pipelineNames, _ := cmd.Flags().GetStringSlice("pipeline")
	jobNames, _ := cmd.Flags().GetStringSlice("job")

	count := 0

	if all {
		count++
	}

	if len(pipelineNames) > 0 {
		count++
	}

	if len(jobNames) > 0 {
		count++
	}

	if count > 1 {
		return errors.New("--all, --pipeline and --job are mutually exclusive")
	}

	if count == 0 {
		return errors.New("either --all, --pipeline or --job must be specified")
	}

	return nil
}
