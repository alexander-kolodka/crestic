package cmd

import (
	"errors"
	"fmt"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// getRepos returns the list of repositories to operate on based on command flags.
// It checks for either --repo flag (specific repositories) or --all flag (all repositories).
// Returns an error if neither flag is specified or if specified repository names are invalid.
func getRepos(cmd *cobra.Command, cfg *entity.Config) ([]*entity.Repository, error) {
	repoNames, _ := cmd.Flags().GetStringSlice("repo")
	err := validateGivenRepoNames(cfg, repoNames)
	if err != nil {
		return nil, err
	}

	all, _ := cmd.Flags().GetBool("all")
	repos := lo.FilterValues(cfg.Repositories, func(_ string, r *entity.Repository) bool {
		return all || lo.Contains(repoNames, r.Name)
	})

	if len(repos) == 0 {
		return nil, errors.New("either --repo or --all must be specified")
	}

	return repos, nil
}

// validateGivenRepoNames checks that all specified repository names exist in the configuration.
// Returns an error if any repository name is not found in the config.
func validateGivenRepoNames(cfg *entity.Config, repoNames []string) error {
	for _, repoName := range repoNames {
		_, ok := cfg.Repositories[repoName]
		if !ok {
			return fmt.Errorf("invalid repository name: %s", repoName)
		}
	}

	return nil
}

// toZerologLevel converts a string log level to a zerolog.Level.
// Supported levels: debug, info, warn, error.
// Returns the corresponding zerolog level or zero value if level is unknown.
func toZerologLevel(level string) zerolog.Level {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
	}

	return levels[level]
}
