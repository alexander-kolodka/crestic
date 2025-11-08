package entity

// Config represents the top-level configuration for crestic.
// It contains all backup/copy jobs, repository definitions, and global settings.
type Config struct {
	Jobs           Jobs                   // List of backup and copy jobs to execute
	Repositories   map[string]*Repository // Map of repository names to repository configs
	HealthcheckURL string                 // Global healthcheck URL for monitoring (can be overridden per job)
}

// Jobs is a list of Job interfaces representing different types of backup operations.
type Jobs []Job

// Job is the interface that all job types (backup, copy) must implement.
// It provides common methods for accessing job properties.
type Job interface {
	GetName() string           // Returns the unique name of the job
	GetHooks() Hooks           // Returns the lifecycle hooks for the job
	GetHealthcheckURL() string // Returns the healthcheck URL for monitoring
	GetCron() string           // Returns the cron expression for scheduling
}

// BackupJob represents a backup operation that backs up directories to a repository.
type BackupJob struct {
	Name                     string      // Unique identifier for this backup job
	HealthcheckURL           string      // Optional healthcheck URL (overrides global setting)
	Cron                     string      // Cron expression for scheduling (e.g., "0 2 * * *")
	IgnoreMissingXAttrsError bool        // If true, ignore extended attributes errors during backup
	From                     []string    // List of source directories to back up
	To                       *Repository // Target repository for storing backups
	Options                  Options     // Additional restic options (tags, excludes, etc.)
	Hooks                    Hooks       // Lifecycle hooks (before, success, failure)
}

// GetName returns the name of the backup job.
func (b BackupJob) GetName() string {
	return b.Name
}

// GetHooks returns the lifecycle hooks configured for this backup job.
func (b BackupJob) GetHooks() Hooks {
	return b.Hooks
}

// GetHealthcheckURL returns the healthcheck URL for monitoring this backup job.
func (b BackupJob) GetHealthcheckURL() string {
	return b.HealthcheckURL
}

// GetCron returns the cron expression for scheduling this backup job.
func (b BackupJob) GetCron() string {
	return b.Cron
}

// CopyJob represents a copy operation that replicates snapshots between repositories.
// This is useful for creating off-site backups or maintaining multiple backup copies.
type CopyJob struct {
	Name           string      // Unique identifier for this copy job
	HealthcheckURL string      // Optional healthcheck URL (overrides global setting)
	Cron           string      // Cron expression for scheduling (e.g., "0 3 * * *")
	From           *Repository // Source repository to copy from
	To             *Repository // Destination repository to copy to
	Options        Options     // Additional restic copy options (tags, filters, etc.)
	Hooks          Hooks       // Lifecycle hooks (before, success, failure)
}

// GetName returns the name of the copy job.
func (c CopyJob) GetName() string {
	return c.Name
}

// GetHooks returns the lifecycle hooks configured for this copy job.
func (c CopyJob) GetHooks() Hooks {
	return c.Hooks
}

// GetHealthcheckURL returns the healthcheck URL for monitoring this copy job.
func (c CopyJob) GetHealthcheckURL() string {
	return c.HealthcheckURL
}

// GetCron returns the cron expression for scheduling this copy job.
// Returns empty string if no schedule is configured.
func (c CopyJob) GetCron() string {
	return c.Cron
}

// Repository represents a restic backup repository configuration.
// It defines where backups are stored and how to access them.
type Repository struct {
	Name          string  // Unique name for this repository
	Path          string  // Repository path or URL (local path, sftp://, s3://, rclone:, etc.)
	PasswordCMD   string  // Shell command that outputs the repository password
	ForgetOptions Options // Retention policy options (keep-daily, keep-weekly, etc.)
}

// Hooks defines lifecycle hooks that run at different stages of a job.
// Hooks are shell commands executed in the specified order.
type Hooks struct {
	Before  []string // Commands to run before the job starts (if any fail, job is aborted)
	Failure []string // Commands to run if the job fails
	Success []string // Commands to run if the job succeeds
}
