package cmd

import (
	"github.com/spf13/cobra"

	"github.com/alexander-kolodka/crestic/internal/cases/cron"
	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	enginecron "github.com/alexander-kolodka/crestic/internal/engine/cron"
	"github.com/alexander-kolodka/crestic/internal/engine/hooks"
	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
	"github.com/alexander-kolodka/crestic/internal/engine/pipelines"
	"github.com/alexander-kolodka/crestic/internal/engine/restic"
	"github.com/alexander-kolodka/crestic/internal/engine/shell"
	"github.com/alexander-kolodka/crestic/internal/pkg/paths"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Run scheduled pipelines based on cron expressions",
	Long: `Run scheduled pipelines based on cron expressions defined in the configuration.

This command is designed to be called periodically by your system scheduler
(e.g., cron, systemd timer, launchd). It intelligently tracks which pipelines are due
to run and executes only those that should run based on their cron schedules.

Key features:
  - State tracking: Remembers last run time to prevent missed or duplicate pipelines
  - File locking: Only one instance can run at a time
  - Flexible scheduling: Can be run every 5, 15, or 30 minutes
  - No missed pipelines: Even if called infrequently, all scheduled pipelines will run

The command:
  1. Loads the last execution state from disk
  2. Checks which pipelines are due to run based on cron expressions
  3. Executes all due pipelines
  4. Saves the current time to state file
  5. Exits (next invocation will start from saved time)

Setup example (add to crontab):
  */5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml --healthcheck

This runs the scheduler every 5 minutes, but pipelines only execute when their
cron expression matches.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		cfg, err := loadConfig(cfgPath)
		if err != nil {
			return err
		}

		resolvedPath, err := findConfigFile(cfgPath)
		if err != nil {
			return err
		}

		canonicalPath, err := paths.Canonical(resolvedPath)
		if err != nil {
			return err
		}

		lockFile := enginecron.LockFileName(canonicalPath)
		stateFile := enginecron.StateFileName(canonicalPath)

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		healthcheck, _ := cmd.Flags().GetBool("healthcheck")

		executor := shell.NewExecutor()
		hooksRunner := hooks.New(executor)
		jobRunner := jobs.NewRunner(restic.NewService(executor), hooksRunner)
		pipelinesRunner := pipelines.NewRunner(jobRunner, hooksRunner)

		h := handler.Chain(
			cron.NewHandler(pipelinesRunner),
			handler.WithPanicRecovery[*cron.Command](),
			handler.WithLock[*cron.Command](lockFile),
		)

		return h.Handle(cmd.Context(), &cron.Command{
			Pipelines:   cfg.Pipelines,
			StateFile:   stateFile,
			DryRun:      dryRun,
			Healthcheck: healthcheck,
		})
	},
}

func init() {
	rootCmd.AddCommand(cronCmd)
	cronCmd.Flags().Bool("dry-run", false, "Dry run")
	cronCmd.Flags().Bool("healthcheck", false, "Send healthcheck pings for pipelines with healthcheck_url")
}
