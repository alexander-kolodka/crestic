package cmd

import (
	"context"
	"errors"

	"github.com/alexander-kolodka/crestic/internal/logger"
	"github.com/alexander-kolodka/crestic/internal/shell"
	"github.com/alexander-kolodka/crestic/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "crestic",
	Version: version.String(),
	Short:   "Configuration-driven backup tool built on restic",
	Long: `Crestic is a wrapper around restic that adds configuration management,
scheduling, health monitoring, and lifecycle hooks to your backup workflow.

Features:
  - YAML-based configuration for all backup jobs
  - Built-in cron scheduler for automated backups
  - Healthchecks.io integration for monitoring
  - Lifecycle hooks (before/success/failure)
  - Easy snapshot replication between repositories
  - Secure password management integration`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		ci, _ := cmd.Flags().GetBool("ci")
		json, _ := cmd.Flags().GetBool("json")

		if ci && json {
			return errors.New("--ci and --json options cannot be used together")
		}

		logLevel, _ := cmd.Flags().GetString("log-level")
		ctx := logger.New(logFormat(ci, json), toZerologLevel(logLevel)).
			WithContext(cmd.Context())
		ctx = logger.WithSource(ctx, "crestic")
		if json {
			ctx = logger.WithJSONMode(ctx)
		}

		cmd.Context()

		printCommands, _ := cmd.Flags().GetBool("print-commands")
		if printCommands {
			ctx = shell.WithPrintingCommands(ctx)
		}

		cmd.SetContext(ctx)

		return nil
	},
}

func Execute(ctx context.Context) error {
	return rootCmd.ExecuteContext(ctx)
}

func init() {
	rootCmd.SetVersionTemplate("crestic version {{.Version}}\n")
	rootCmd.PersistentFlags().
		StringP("config", "c", "", "config file (default is crestic.yml in current dir,"+
			" home dir, ~/.crestic/, or ~/.config/crestic/)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().Bool("ci", false, "output logs as plain text without colors (for CI/pipelines)")
	rootCmd.PersistentFlags().Bool("json", false, "output logs in JSON format")
	rootCmd.PersistentFlags().Bool("print-commands", false, "Print executed shell commands")

	_ = rootCmd.MarkFlagFilename("config", "yaml", "yml")

	_ = rootCmd.RegisterFlagCompletionFunc(
		"log-level",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"debug", "info", "warn", "error"}, cobra.ShellCompDirectiveDefault
		},
	)
}

func logFormat(ci, json bool) logger.Format {
	if ci {
		return logger.FormatCI
	}

	if json {
		return logger.FormatJSON
	}

	return logger.FormatColor
}
