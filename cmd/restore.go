package cmd

import (
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/cases/restore"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a snapshot to a directory",
	Long: `Restore a backup snapshot from a repository to a specified directory.

This command extracts files from a backup snapshot to a target directory.
You can restore the entire snapshot or specific files/directories. By default,
the latest snapshot is restored.

Examples:
  # Restore latest snapshot to a directory
  crestic restore --repo local-backup --target ./restore

  # Restore specific snapshot by ID
  crestic restore --repo local-backup --snapshot abc123 --target ./restore

First list snapshots to see what's available:
  crestic exec --repo local-backup snapshots`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		repo, _ := cmd.Flags().GetString("repo")
		err = validateGivenRepoNames(cfg, []string{repo})
		if err != nil {
			return err
		}

		target, _ := cmd.Flags().GetString("target")
		snapshot, _ := cmd.Flags().GetString("snapshot")
		if snapshot == "" {
			snapshot = "latest"
		}

		executor := shell.NewExecutor()
		h := handler.Chain(
			restore.NewHandler(restic.NewService(executor)),
			handler.WithPanicRecovery[*restore.Command](),
		)

		return h.Handle(cmd.Context(), &restore.Command{
			Repo:     cfg.Repositories[repo],
			Target:   target,
			Snapshot: snapshot,
		})
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().StringP("repo", "r", "", "Restore specific repository (required)")
	restoreCmd.Flags().StringP("target", "t", "", "Directory to restore to (required)")
	restoreCmd.Flags().StringP("snapshot", "s", "", "snapshot")
	_ = restoreCmd.MarkFlagRequired("target")

	_ = restoreCmd.RegisterFlagCompletionFunc("repo", repoAutocompletion)
}
