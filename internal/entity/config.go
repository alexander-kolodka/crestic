package entity

import (
	"fmt"

	"github.com/samber/lo"
)

// Config represents the top-level configuration for crestic.
type Config struct {
	Pipelines    Pipelines
	Repositories map[string]*Repository // Map of repository names to repository configs
}

type Pipelines []Pipeline

// Pipeline groups related jobs that run together on a schedule.
type Pipeline struct {
	Name           string
	Cron           string
	HealthcheckURL string // Optional Healthchecks.io ping URL for this pipeline
	Hooks          Hooks  // Optional lifecycle hooks for the pipeline
	Jobs           Jobs
}

// Jobs is a list of Job interfaces representing different types of backup operations.
type Jobs []Job

func (c *Config) AllJobs() []Job {
	return lo.FlatMap(c.Pipelines, func(p Pipeline, _ int) []Job {
		return p.Jobs
	})
}

func (c *Config) FindJob(fullName string) (Job, error) {
	pipeline, job, err := splitFullName(fullName)
	if err != nil {
		return nil, err
	}

	p, ok := c.FindPipeline(pipeline)
	if !ok {
		return nil, fmt.Errorf("pipeline %s not found", pipeline)
	}

	j, ok := lo.Find(p.Jobs, func(j Job) bool {
		return j.GetName() == job
	})

	if !ok {
		return nil, fmt.Errorf("job %s not found", job)
	}

	return j, nil
}

func (c *Config) FindPipeline(pipeline string) (Pipeline, bool) {
	return lo.Find(c.Pipelines, func(p Pipeline) bool {
		return p.Name == pipeline
	})
}
