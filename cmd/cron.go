package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alexander-kolodka/crestic/internal/cases/backup"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/cron"
	"github.com/alexander-kolodka/crestic/internal/healthchecks"
	"github.com/alexander-kolodka/crestic/internal/restic"
	"github.com/alexander-kolodka/crestic/internal/shell"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Run scheduled jobs based on cron expressions",
	Long: `Run scheduled backup jobs based on cron expressions defined in the configuration.

This command is designed to be called periodically by your system scheduler
(e.g., cron, systemd timer, launchd). It intelligently tracks which jobs are due
to run and executes only those that should run based on their cron schedules.

Key features:
  - State tracking: Remembers last run time to prevent missed or duplicate jobs
  - File locking: Only one instance can run at a time
  - Flexible scheduling: Can be run every 5, 15, or 30 minutes
  - No missed jobs: Even if called infrequently, all scheduled jobs will run

The command:
  1. Loads the last execution state from disk
  2. Checks which jobs are due to run based on cron expressions
  3. Executes all due jobs
  4. Saves the current time to state file
  5. Exits (next invocation will start from saved time)

Setup example (add to crontab):
  */5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml

This runs the scheduler every 5 minutes, but jobs only execute when their
cron expression matches.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		fileName, err := getCfgFileName(cfgPath)
		if err != nil {
			return err
		}

		jobs, err := cron.FilterJobsByCron(cmd.Context(), cfg.Jobs)
		if err != nil {
			return err
		}

		executor := shell.NewExecutor()
		hcClient := healthchecks.NewClient()
		h := handler.Chain(
			backup.NewHandler(restic.NewService(executor), executor, hcClient),
			handler.WithPanicRecovery[*backup.Command](),
			handler.WithLock[*backup.Command](fmt.Sprintf("crestic-cron-%s.lock", fileName)),
		)

		return h.Handle(cmd.Context(), &backup.Command{
			Jobs: jobs,
		})
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
}

func getCfgFileName(cfgPath string) (string, error) {
	path, err := findConfigFile(cfgPath)
	if err != nil {
		return "", err
	}

	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return fileName, nil
}
