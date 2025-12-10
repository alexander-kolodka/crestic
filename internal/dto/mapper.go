package dto

import (
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/alexander-kolodka/crestic/internal/entity"
)

func ToEntity(cfg Config) (*entity.Config, error) {
	repos := lo.MapEntries(cfg.Repositories,
		func(name string, repo Repository) (string, *entity.Repository) {
			return name, toRepository(name, repo)
		},
	)

	missedRepos := make(map[string]struct{})

	jobs := lo.Map(cfg.Jobs,
		func(job Job, _ int) entity.Job {
			switch j := job.(type) {
			case BackupJob:
				repo, ok := repos[j.To]
				if !ok {
					missedRepos[j.To] = struct{}{}
				}

				return toBackupJob(j, repo)
			case CopyJob:
				from, ok := repos[j.From]
				if !ok {
					missedRepos[j.From] = struct{}{}
				}

				to, ok := repos[j.To]
				if !ok {
					missedRepos[j.To] = struct{}{}
				}

				return toCopyJob(j, from, to)
			default:
			}

			return nil
		},
	)

	missed := lo.Keys(missedRepos)
	if len(missed) > 0 {
		return nil, fmt.Errorf("missed repositories: %s", strings.Join(missed, ", "))
	}

	return &entity.Config{
		HealthcheckURL: cfg.HealthcheckURL,
		Repositories:   repos,
		Jobs:           jobs,
	}, nil
}

func toRepository(name string, repo Repository) *entity.Repository {
	return &entity.Repository{
		Name:          name,
		Path:          repo.Path,
		PasswordCMD:   repo.PasswordCMD,
		ForgetOptions: entity.Options(repo.ForgetOptions),
	}
}

func toBackupJob(b BackupJob, repo *entity.Repository) entity.BackupJob {
	return entity.BackupJob{
		Name:                     b.Name,
		Cron:                     b.Cron,
		IgnoreMissingXAttrsError: b.IgnoreMissingXAttrsError,
		From:                     b.From,
		To:                       repo,
		Options:                  entity.Options(b.Options),
		Hooks:                    toHooks(b.Hooks),
	}
}

func toCopyJob(c CopyJob, from, to *entity.Repository) entity.CopyJob {
	return entity.CopyJob{
		Name:    c.Name,
		Cron:    c.Cron,
		From:    from,
		To:      to,
		Options: entity.Options(c.Options),
		Hooks:   toHooks(c.Hooks),
	}
}

func toHooks(h Hooks) entity.Hooks {
	return entity.Hooks{
		Before:  h.Before,
		Failure: h.Failure,
		Success: h.Success,
	}
}
