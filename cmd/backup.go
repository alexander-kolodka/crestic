package cmd

import (
	"errors"

	"github.com/alexander-kolodka/crestic/internal/cases/backup"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create encrypted backups of configured sources",
	Long: `Create encrypted backups of directories specified in your configuration file.

This command performs a complete backup workflow for each job:

The backup process (automatic):
  1. Runs 'before' hooks (if configured)
  2. Sends start notification to healthcheck service (if configured)
  3. Checks if repository is initialized (auto-initializes if needed)
  4. Creates encrypted backup snapshot using restic
  5. Verifies repository integrity (restic check)
  6. Applies retention policy (restic forget with forget_options)
  7. Sends success/failure notification to healthcheck service
  8. Runs 'success' or 'failure' hooks based on outcome

Note: The forget step automatically runs after each backup if forget_options
are configured in the repository. If --prune flag is set in forget_options,
old data is actually removed from the repository to free disk space.

A failure in one backup job doesn't prevent other backups from completing.
At the end, all errors are collected and returned as a combined error.

Examples:
  # Backup all configured jobs
  crestic backup --all

  # Backup specific job
  crestic backup --job documents

  # Backup multiple jobs
  crestic backup --job documents,photos

  # Dry run (show what would be backed up)
  crestic backup --all --dry-run`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		jobs := filterJobs(cmd, cfg.Jobs)
		if len(jobs) == 0 {
			return errors.New("either --job or --all must be specified")
		}

		executor := shell.NewExecutor()
		hcClient := healthchecks.NewClient()
		h := handler.Chain(
			backup.NewHandler(restic.NewService(executor), executor, hcClient),
			handler.WithPanicRecovery[*backup.Command](),
		)

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		return h.Handle(cmd.Context(), &backup.Command{
			Jobs:   jobs,
			DryRun: dryRun,
		})
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().BoolP("all", "a", false, "Check all repositories")
	backupCmd.Flags().StringSliceP("job", "j", nil, "Run only specific jobs by name (comma-separated)")
	backupCmd.Flags().Bool("dry-run", false, "Dry run")

	_ = backupCmd.RegisterFlagCompletionFunc("job", jobAutocompletion)
}

func filterJobs(cmd *cobra.Command, backups []entity.Job) []entity.Job {
	all, _ := cmd.Flags().GetBool("all")
	jNames, _ := cmd.Flags().GetStringSlice("job")
	return lo.Filter(backups, func(b entity.Job, _ int) bool {
		return all || lo.Contains(jNames, b.GetName())
	})
}
