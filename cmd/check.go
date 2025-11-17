package cmd

import (
	"github.com/spf13/cobra"

	"github.com/alexander-kolodka/crestic/internal/cases/check"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check and initialize repositories",
	Long: `Check that all configured repositories are properly initialized and accessible.

This command verifies each repository's existence and integrity. If a repository
doesn't exist, it will be initialized automatically. This is typically the first
command you run after creating a new configuration.

For each repository, the command:
  1. Checks if the repository is initialized
  2. Creates a new repository if it doesn't exist
  3. Verifies repository integrity if it does exist
  4. Reports any errors or issues found

Examples:
  # Check all repositories
  crestic check --all

  # Check specific repository
  crestic check --repo local-backup

  # Check multiple repositories
  crestic check --repo local-backup,remote-backup`,
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
			check.NewHandler(restic.NewService(executor)),
			handler.WithPanicRecovery[*check.Command](),
		)

		return h.Handle(cmd.Context(), &check.Command{
			Repos: repos,
		})
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringSliceP("repo", "r", nil, "Check specific repository/repositories (can specify multiple)")
	checkCmd.Flags().BoolP("all", "a", false, "Check all repositories")
}
