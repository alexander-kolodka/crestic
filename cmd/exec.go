package cmd

import (
	"errors"

	"github.com/alexander-kolodka/crestic/internal/cases/exec"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [flags] <command> [-- restic-options]",
	Short: "Execute native restic commands on repositories",
	Long: `Execute any native restic command on configured repositories.

This command provides direct access to restic functionality while automatically
using the repository paths and passwords from your crestic configuration. This
is useful for advanced operations not directly supported by crestic commands.

The command:
  1. Loads repository configuration from crestic.yaml
  2. Sets up repository path and password for each repository
  3. Executes the specified restic command
  4. Returns the results from each repository

Common restic commands:
  snapshots   - List all snapshots
  ls          - List files in a snapshot
  find        - Find a file across all snapshots
  stats       - Show repository statistics
  diff        - Compare two snapshots
  mount       - Mount repository as filesystem (requires FUSE)
  check       - Perform deep integrity check
  prune       - Remove unreferenced data
  rebuild-index - Rebuild repository index

Examples:
  # List all snapshots for all repositories
  crestic exec --all snapshots

  # List snapshots for specific repository
  crestic exec --repo local-backup snapshots

  # List files in latest snapshot
  crestic exec --repo local-backup ls latest

  # Find a file in all snapshots
  crestic exec --repo local-backup find document.pdf

  # Show repository statistics
  crestic exec --repo local-backup stats

  # Compare two snapshots
  crestic exec --repo local-backup diff abc123 def456

  # Mount repository (requires FUSE)
  crestic exec --repo local-backup mount /mnt/backup

For more information about restic commands, see:
  https://restic.readthedocs.io/`,
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("no command specified")
		}

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
			exec.NewHandler(restic.NewService(executor)),
			handler.WithPanicRecovery[*exec.Command](),
		)

		return h.Handle(cmd.Context(), &exec.Command{
			Repos: repos,
			Cmd:   args[0],
			Args:  args[1:],
		})
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().
		StringSliceP("repo", "r", nil, "Repository/repositories to execute command on (can specify multiple)")
	execCmd.Flags().BoolP("all", "a", false, "Execute on all repositories")
}
