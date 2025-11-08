package cron

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const stateFileName = "crestic-cron-state.json"

type State struct {
	LastRun time.Time `json:"last_run"`
}

// loadState loads the last run time from the state file.
func loadState() (time.Time, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load home dir: %w", err)
	}

	statePath := filepath.Join(homeDir, ".crestic", stateFileName)
	data, err := os.ReadFile(statePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to read state file: %w", err)
	}

	var state State
	err = json.Unmarshal(data, &state)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return state.LastRun, nil
}

// saveState saves the current time as the last run time.
func saveState(now time.Time) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	cresticDir := filepath.Join(homeDir, ".crestic")
	err = os.MkdirAll(cresticDir, 0o750)
	if err != nil {
		return fmt.Errorf("failed to create .crestic directory: %w", err)
	}

	statePath := filepath.Join(cresticDir, stateFileName)

	state := State{
		LastRun: now,
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
