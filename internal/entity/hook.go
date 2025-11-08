package entity

import "time"

type HookStart struct {
	JobName string `json:"JobName"`
}

type HookSuccess struct {
	JobName     string         `json:"JobName"`
	Elapsed     time.Duration  `json:"Elapsed"`
	BackupStats map[string]any `json:"BackupStats"`
	CopyStats   map[string]any `json:"CopyStats"`
	ForgetStats map[string]any `json:"ForgetStats"`
}

type HookFailure struct {
	JobName  string        `json:"JobName"`
	Elapsed  time.Duration `json:"Elapsed"`
	Error    string        `json:"Error"`
	ExitCode int           `json:"ExitCode"`
}
