package dto

import (
	"errors"
	"fmt"
	"strings"

	"github.com/samber/lo"

	"github.com/alexander-kolodka/crestic/internal/entity"
)

func ToEntity(cfg Config) (*entity.Config, error) {
	repos := toRepositories(cfg.Repositories)
	lookup := newRepoLookup(repos)

	pipelines, err := toPipelines(cfg.Pipelines, lookup)
	if err != nil {
		return nil, err
	}

	err = lookup.err()
	if err != nil {
		return nil, err
	}

	return &entity.Config{
		Repositories: repos,
		Pipelines:    pipelines,
	}, nil
}

type repoLookup struct {
	repos  map[string]*entity.Repository
	missed map[string]struct{}
}

func newRepoLookup(repos map[string]*entity.Repository) *repoLookup {
	return &repoLookup{repos: repos, missed: make(map[string]struct{})}
}

func (r *repoLookup) get(name string) *entity.Repository {
	repo, ok := r.repos[name]
	if !ok {
		r.missed[name] = struct{}{}
	}
	return repo
}

func (r *repoLookup) err() error {
	if len(r.missed) == 0 {
		return nil
	}
	return fmt.Errorf("missed repositories: %s", strings.Join(lo.Keys(r.missed), ", "))
}

func toRepositories(repos map[string]Repository) map[string]*entity.Repository {
	return lo.MapEntries(repos,
		func(name string, repo Repository) (string, *entity.Repository) {
			return name, toRepository(name, repo)
		},
	)
}

func toPipelines(pipelines Pipelines, repos *repoLookup) (entity.Pipelines, error) {
	seen := make(map[string]struct{}, len(pipelines))
	out := make(entity.Pipelines, 0, len(pipelines))

	for _, p := range pipelines {
		if p.Name == "" {
			return nil, errors.New("pipeline name is required")
		}
		if _, exists := seen[p.Name]; exists {
			return nil, fmt.Errorf("duplicate pipeline name: %s", p.Name)
		}
		seen[p.Name] = struct{}{}

		ep, err := toPipeline(p, repos)
		if err != nil {
			return nil, err
		}
		out = append(out, ep)
	}

	return out, nil
}

func toPipeline(p Pipeline, repos *repoLookup) (entity.Pipeline, error) {
	jobs, err := toJobs(p.Jobs, p.Name, repos)
	if err != nil {
		return entity.Pipeline{}, err
	}
	return entity.Pipeline{
		Name:           p.Name,
		Cron:           p.Cron,
		HealthcheckURL: p.HealthcheckURL,
		Hooks:          toHooks(p.Hooks),
		Jobs:           jobs,
	}, nil
}

func toJobs(jobs Jobs, pipeline string, repos *repoLookup) (entity.Jobs, error) {
	seen := make(map[string]struct{}, len(jobs))
	out := make(entity.Jobs, 0, len(jobs))

	for _, job := range jobs {
		entityJob, err := toJob(job, pipeline, repos)
		if err != nil {
			return nil, err
		}
		name := entityJob.GetName()
		if _, exists := seen[name]; exists {
			return nil, fmt.Errorf("duplicate job name %q in pipeline %q", name, pipeline)
		}
		seen[name] = struct{}{}
		out = append(out, entityJob)
	}

	return out, nil
}

func toJob(job Job, pipeline string, repos *repoLookup) (entity.Job, error) {
	switch j := job.(type) {
	case BackupJob:
		return toBackupJob(j, pipeline, repos.get(j.To)), nil
	case CopyJob:
		return toCopyJob(j, pipeline, repos.get(j.From), repos.get(j.To)), nil
	default:
		return nil, fmt.Errorf("unknown job type in pipeline %q", pipeline)
	}
}

func toRepository(name string, repo Repository) *entity.Repository {
	return &entity.Repository{
		Name:          name,
		Path:          repo.Path,
		PasswordCMD:   repo.PasswordCMD,
		ForgetOptions: entity.Options(repo.ForgetOptions),
	}
}

func toBackupJob(b BackupJob, pipeline string, repo *entity.Repository) entity.BackupJob {
	return entity.BackupJob{
		Name:                     b.Name,
		Pipeline:                 pipeline,
		IgnoreMissingXAttrsError: b.IgnoreMissingXAttrsError,
		From:                     b.From,
		To:                       repo,
		Options:                  entity.Options(b.Options),
		Hooks:                    toHooks(b.Hooks),
	}
}

func toCopyJob(c CopyJob, pipeline string, from, to *entity.Repository) entity.CopyJob {
	return entity.CopyJob{
		Name:     c.Name,
		Pipeline: pipeline,
		From:     from,
		To:       to,
		Options:  entity.Options(c.Options),
		Hooks:    toHooks(c.Hooks),
	}
}

func toHooks(h Hooks) entity.Hooks {
	return entity.Hooks{
		Before:  h.Before,
		Failure: h.Failure,
		Success: h.Success,
	}
}
