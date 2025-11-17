package cmd

import (
	"github.com/spf13/cobra"

	"github.com/alexander-kolodka/crestic/internal/cases/forget"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

var forgetCmd = &cobra.Command{
	Use:   "forget",
	Short: "Remove old snapshots according to retention policy",
	Long: `Remove old backup snapshots according to the retention policy defined in configuration.

This command applies the forget_options configured for each repository to remove
old snapshots that are no longer needed. This is essential for managing disk space
and maintaining a reasonable number of backups.

Understanding forget vs prune:
  - forget: Marks snapshots for deletion (fast, metadata only)
  - prune: Actually removes data from repository (slow, frees disk space)

Retention policy example:
  forget_options:
    keep-daily: 7      # Keep 7 daily snapshots
    keep-weekly: 4     # Keep 4 weekly snapshots
    keep-monthly: 12   # Keep 12 monthly snapshots
    keep-yearly: 3     # Keep 3 yearly snapshots

The command will keep the specified number of snapshots for each time period
and mark older ones for deletion.

Examples:
  # Show what would be deleted (safe, no changes)
  crestic forget --all --dry-run

  # Mark old snapshots for deletion (fast, doesn't free space)
  crestic forget --all

  # Mark old snapshots and remove data (slow, frees disk space)
  crestic forget --all --prune

  # Forget for specific repository
  crestic forget --repo local-backup --prune`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		repos, err := getRepos(cmd, cfg)
		if err != nil {
			return err
		}

		executor := shell.NewExecutor()
		h := handler.Chain(
			forget.NewHandler(restic.NewService(executor)),
			handler.WithPanicRecovery[*forget.Command](),
		)

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		prune, _ := cmd.Flags().GetBool("prune")
		return h.Handle(cmd.Context(), &forget.Command{
			Repos:  repos,
			Prune:  prune,
			DryRun: dryRun,
		})
	},
}

func init() {
	rootCmd.AddCommand(forgetCmd)
	forgetCmd.Flags().StringSliceP("repo", "r", nil, "Run forget for a specific repository")
	forgetCmd.Flags().BoolP("all", "a", false, "Run forget for all jobs")
	forgetCmd.Flags().Bool("dry-run", false, "Show what would be deleted without actually deleting")
	forgetCmd.Flags().Bool("prune", false, "Actually remove the data (frees up space)")

	_ = forgetCmd.RegisterFlagCompletionFunc("repo", repoAutocompletion)
}
