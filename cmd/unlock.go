package cmd

import (
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/cases/unlock"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Remove stale locks from repositories",
	Long: `Remove stale locks from repositories that were not automatically released.

Restic uses locks to prevent concurrent access to repositories. Normally, locks
are automatically released when operations complete. However, locks may remain
if a process was killed, crashed, or lost network connectivity.

When to use this command:
  - You see "repository is already locked" errors
  - You're certain no other backup operation is running
  - A previous operation was interrupted (crash, kill, network loss)

WARNING: Only run this command if you're absolutely certain that no other
crestic or restic process is currently accessing the repository. Running unlock
while another operation is in progress can cause data corruption.

Examples:
  # Unlock all repositories
  crestic unlock --all

  # Unlock specific repository
  crestic unlock --repo local-backup

  # Unlock multiple repositories
  crestic unlock --repo local-backup,remote-backup`,
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
			unlock.NewHandler(restic.NewService(executor)),
			handler.WithPanicRecovery[*unlock.Command](),
		)

		return h.Handle(cmd.Context(), &unlock.Command{
			Repos: repos,
		})
	},
}

func init() {
	rootCmd.AddCommand(unlockCmd)
	unlockCmd.Flags().StringSliceP("repo", "r", nil, "Unlock specific repository/repositories (can specify multiple)")
	unlockCmd.Flags().BoolP("all", "a", false, "Unlock all repositories")

	_ = unlockCmd.RegisterFlagCompletionFunc("repo", repoAutocompletion)
}
