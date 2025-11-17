package cron

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/samber/lo"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

// FilterJobsByCron filters jobs that should run based on their cron expressions.
// It automatically loads the last run time from state, checks all jobs that were scheduled
// to run since then, and saves the current time as the new last run time.
// If state file doesn't exist (first run), it uses current time as lastRun to avoid running all historical jobs.
func FilterJobsByCron(ctx context.Context, jobs []entity.Job) ([]entity.Job, error) {
	log := logger.FromContext(ctx)

	now := time.Now()

	lastRun, err := loadState()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load state, using current time as last run")
		lastRun = now
	}

	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

	jobs = lo.Filter(jobs, func(job entity.Job, _ int) bool {
		if job.GetCron() == "" {
			log.Debug().
				Str("job", job.GetName()).
				Msg("Skip job with no cron expression")
			return false
		}

		schedule, parseErr := cronParser.Parse(job.GetCron())
		if parseErr != nil {
			log.Warn().Err(parseErr).
				Str("job", job.GetName()).
				Str("cron", job.GetCron()).
				Msg("Failed to parse cron expression, skipping job")
			return false
		}

		runAt := schedule.Next(lastRun)
		if !runAt.Before(now) {
			log.Debug().
				Str("job", job.GetName()).
				Str("cron", job.GetCron()).
				Time("run_at", runAt).
				Time("now", now).
				Msg("Skip job")
			return false
		}

		log.Debug().
			Str("job", job.GetName()).
			Str("cron", job.GetCron()).
			Time("run_at", runAt).
			Time("now", now).
			Msg("Process job")

		return true
	})

	err = saveState(now)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save state")
		return nil, err
	}

	return jobs, nil
}
