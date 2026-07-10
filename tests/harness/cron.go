package harness

import (
	"time"
)

const cronStateFileName = "crestic-cron-state.json"

type cronState struct {
	Pipelines map[string]cronPipelineState `json:"pipelines"`
}

type cronPipelineState struct {
	LastRun time.Time `json:"last_run"`
}
