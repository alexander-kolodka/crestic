package cron

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Pipelines map[string]PipelineState `json:"pipelines"`
}

type PipelineState struct {
	LastRun time.Time `json:"last_run"`
}

func loadState(stateFile string) (State, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return State{}, fmt.Errorf("failed to get home directory: %w", err)
	}

	statePath := filepath.Join(homeDir, ".crestic", stateFile)
	data, err := os.ReadFile(statePath)
	if err != nil {
		return State{}, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	err = json.Unmarshal(data, &state)
	if err != nil {
		return State{}, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	if state.Pipelines == nil {
		state.Pipelines = map[string]PipelineState{}
	}

	return state, nil
}

func saveState(stateFile string, state State) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	cresticDir := filepath.Join(homeDir, ".crestic")
	err = os.MkdirAll(cresticDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create .crestic directory: %w", err)
	}

	statePath := filepath.Join(cresticDir, stateFile)

	if state.Pipelines == nil {
		state.Pipelines = map[string]PipelineState{}
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	err = os.WriteFile(statePath, data, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}
