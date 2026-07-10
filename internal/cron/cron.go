package cron

import (
	"context"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/alexander-kolodka/crestic/internal/entity"
	"github.com/alexander-kolodka/crestic/internal/logger"
)

// FilterPipelinesByCron filters pipelines that should run based on their cron expressions.
// It loads per-pipeline last run times from state, checks which pipelines were scheduled
// to run since then, and saves the current time for each due pipeline.
func FilterPipelinesByCron(ctx context.Context, pipelines entity.Pipelines) (entity.Pipelines, error) {
	log := logger.FromContext(ctx)

	now := time.Now()

	state, err := loadState()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load state, using empty state")
		state = State{Pipelines: map[string]PipelineState{}}
	}

	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

	var duePipelines entity.Pipelines
	stateChanged := false

	for _, pipeline := range pipelines {
		if pipeline.Cron == "" {
			log.Debug().
				Str("pipeline", pipeline.Name).
				Msg("Skip pipeline with no cron expression")
			continue
		}

		schedule, parseErr := cronParser.Parse(pipeline.Cron)
		if parseErr != nil {
			log.Warn().Err(parseErr).
				Str("pipeline", pipeline.Name).
				Str("cron", pipeline.Cron).
				Msg("Failed to parse cron expression, skipping pipeline")
			continue
		}

		lastRun, ok := state.Pipelines[pipeline.Name]
		if !ok {
			lastRun = PipelineState{LastRun: now}
		}

		runAt := schedule.Next(lastRun.LastRun)
		if !runAt.Before(now) {
			log.Debug().
				Str("pipeline", pipeline.Name).
				Str("cron", pipeline.Cron).
				Time("run_at", runAt).
				Time("now", now).
				Msg("Skip pipeline")
			continue
		}

		log.Debug().
			Str("pipeline", pipeline.Name).
			Str("cron", pipeline.Cron).
			Time("run_at", runAt).
			Time("now", now).
			Msg("Process pipeline")

		duePipelines = append(duePipelines, pipeline)
		state.Pipelines[pipeline.Name] = PipelineState{LastRun: now}
		stateChanged = true
	}

	if stateChanged {
		err = saveState(state)
		if err != nil {
			log.Error().Err(err).Msg("Failed to save state")
			return nil, err
		}
	}

	return duePipelines, nil
}
