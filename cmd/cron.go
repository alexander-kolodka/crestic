package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alexander-kolodka/crestic/internal/cases/handler"
	"github.com/alexander-kolodka/crestic/internal/cases/runpipelines"
	"github.com/alexander-kolodka/crestic/internal/engine/cron"
	"github.com/alexander-kolodka/crestic/internal/engine/hooks"
	"github.com/alexander-kolodka/crestic/internal/engine/jobs"
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
  */5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml --healthcheck

This runs the scheduler every 5 minutes, but jobs only execute when their
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

		canonicalPath, err := cron.CanonicalConfigPath(resolvedPath)
		if err != nil {
			return err
		}

		cfgName := strings.TrimSuffix(filepath.Base(canonicalPath), filepath.Ext(canonicalPath))
		lockFile := cron.LockFileName(canonicalPath, cfgName)
		stateFile := cron.StateFileName(canonicalPath, cfgName)

		duePipelines, err := cron.FilterPipelinesByCron(cmd.Context(), cfg.Pipelines, stateFile)
		if err != nil {
			return err
		}

		if len(duePipelines) == 0 {
			return nil
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		healthcheck, _ := cmd.Flags().GetBool("healthcheck")

		executor := shell.NewExecutor()
		hooksRunner := hooks.New(executor)
		jobRunner := jobs.NewRunner(restic.NewService(executor), hooksRunner)

		h := handler.Chain(
			runpipelines.NewHandler(jobRunner, hooksRunner),
			handler.WithPanicRecovery[*runpipelines.Command](),
			handler.WithLock[*runpipelines.Command](lockFile),
		)

		return h.Handle(cmd.Context(), &runpipelines.Command{
			Pipelines:   duePipelines,
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
