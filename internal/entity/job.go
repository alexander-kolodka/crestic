package entity

import (
	"fmt"
	"strings"
)

// Job is the interface that all job types (backup, copy) must implement.
type Job interface {
	GetName() string
	GetPipeline() string
	GetFullName() string
	GetHooks() Hooks
}

// BackupJob represents a backup operation that backs up directories to a repository.
type BackupJob struct {
	Name                     string
	Pipeline                 string
	IgnoreMissingXAttrsError bool
	From                     []string
	To                       *Repository
	Options                  Options
	Hooks                    Hooks
}

// GetName returns the name of the backup job.
func (b BackupJob) GetName() string {
	return b.Name
}

func (b BackupJob) GetPipeline() string {
	return b.Pipeline
}

func (b BackupJob) GetFullName() string {
	return b.Pipeline + "/" + b.Name
}

func (b BackupJob) GetHooks() Hooks {
	return b.Hooks
}

// CopyJob represents a copy operation that replicates snapshots between repositories.
type CopyJob struct {
	Name     string
	Pipeline string
	From     *Repository
	To       *Repository
	Options  Options
	Hooks    Hooks
}

func (c CopyJob) GetName() string {
	return c.Name
}

func (c CopyJob) GetPipeline() string {
	return c.Pipeline
}

func (c CopyJob) GetFullName() string {
	return c.Pipeline + "/" + c.Name
}

func (c CopyJob) GetHooks() Hooks {
	return c.Hooks
}

// Repository represents a restic backup repository configuration.
type Repository struct {
	Name          string  // Unique name for this repository
	Path          string  // Repository path or URL (local path, sftp://, s3://, rclone:, etc.)
	PasswordCMD   string  // Shell command that outputs the repository password
	ForgetOptions Options // Retention policy options (keep-daily, keep-weekly, etc.)
}

// Hooks defines lifecycle hooks that run at different stages of a job.
type Hooks struct {
	Before  []string // Commands to run before the job starts (if any fail, job is aborted)
	Failure []string // Commands to run if the job fails
	Success []string // Commands to run if the job succeeds
}

func splitFullName(fullName string) (pipeline, job string, err error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid job name %s", fullName)
	}

	pipeline, job = parts[0], parts[1]
	return pipeline, job, nil
}
