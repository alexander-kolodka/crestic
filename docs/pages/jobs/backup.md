# 💾 Backup Job

Backup jobs back up local directories to a restic repository.

## Configuration Structure

```yaml
jobs:
  - type: backup
    name: string                    # Required: Unique job name
    from: []string                  # Required: Source directories
    to: string                      # Required: Target repository name
    cron: string                    # Optional: Cron expression
    healthcheck_url: string         # Optional: Job-specific healthcheck URL
    ignore_x_attrs_error: bool      # Optional: Ignore extended attributes errors
    options:                        # Optional: Restic backup options
      key: value
    hooks:                          # Optional: Lifecycle hooks
      before: []string
      success: []string
      failure: []string
```

## Required Fields

### `type`

Must be `"backup"`.

### `name`

Unique identifier for the job. Used in logs and when selecting specific jobs.

```yaml
name: documents-backup
name: photos-daily
name: system-files
```

### `from`

List of directories to back up. All paths are included in a single snapshot.

```yaml
from:
  - /home/user/Documents
  - /home/user/Projects
  - /home/user/.config
```

**Note**: All directories listed in `from` are backed up together in one snapshot.

### `to`

Name of the target repository (must be defined in `repositories` section).

```yaml
to: local-repo
to: remote-backup
```

## Optional Fields

### `cron`

Cron expression for scheduling automated backups.

**Format**: `minute hour day month weekday`

**Examples**:
```yaml
cron: "0 2 * * *"      # Daily at 2:00 AM
cron: "0 */6 * * *"    # Every 6 hours
cron: "30 3 * * 0"     # Weekly on Sunday at 3:30 AM
cron: "0 4 1 * *"      # Monthly on 1st at 4:00 AM
cron: "0 9,17 * * 1-5" # Weekdays at 9 AM and 5 PM
```

**To use scheduling**:
1. Set cron expression in job configuration
2. Add to system crontab: `*/5 * * * * crestic cron --config /path/to/crestic.yaml`
3. Crestic tracks state, so system cron can run every 5-30 minutes

See [Cron Command](/cli/cron) for more details.

### `healthcheck_url`

Override global healthcheck URL for this specific job.

```yaml
healthcheck_url: https://hc-ping.com/uuid-for-documents/documents-backup
```

See [Healthchecks](/healthchecks) for more details.

### `ignore_x_attrs_error`
Some filesystems (e.g. Cryptomator, other FUSE mounts) do not allow reading extended file attributes (xattrs).
When restic encounters such files, it exits with status code 3, which means:

>“incomplete metadata for ${file}”

However, the backup is still created successfully and only the unreadable xattrs are skipped.
If you want Crestic to ignore this exit code and treat the backup as successful, enable the option:

```yaml
ignore_x_attrs_error: true
```

**Default**: `false`

## Options

The `options` field accepts any restic backup option. Common options:

### Tagging

```yaml
options:
  tag:
    - documents
    - daily
    - important
```

### Exclude Patterns

```yaml
options:
  exclude:
    - "*.tmp"
    - "*.log"
    - ".cache"
    - "node_modules"
  exclude-file: "/path/to/exclude.txt"
```

### Include Patterns

```yaml
options:
  files-from: "/path/to/include.txt"
```

### Performance Options

```yaml
options:
  skip-if-unchanged: true  # Skip backup if no files changed
  one-file-system: true    # Don't cross filesystem boundaries
  with-atime: false        # Don't save access time
```

### Other Options

```yaml
options:
  host: "my-server"        # Set hostname for snapshot
  exclude-caches: true     # Exclude cache directories
  exclude-if-present: ".nobackup"  # Exclude if file present
```

**See**: [Restic Backup Documentation](https://restic.readthedocs.io/en/stable/040_backup.html) for complete list of options.

## Hooks

Execute custom commands at different stages:

```yaml
hooks:
  before:
    - echo "Starting backup..."
    - /usr/local/bin/pre-backup-script.sh
  success:
    - echo "Backup completed!"
    - curl -X POST https://your-webhook.com/success
  failure:
    - echo "Backup failed!" >&2
    - /usr/local/bin/alert-admin.sh
```

**Environment variables available in hooks**:
- `CRESTIC_JOB_NAME` - Name of the job
- `CRESTIC_EXIT_CODE` - Exit code of the operation
- `CRESTIC_ERROR` - Error message (only in failure hooks)

See [Hooks](/hooks) for more details.

## Complete Example

```yaml
jobs:
  - type: backup
    name: documents-backup
    from:
      - /home/user/Documents
      - /home/user/Projects
    to: local-repo
    cron: "0 2 * * *"
    healthcheck_url: https://hc-ping.com/uuid/documents
    ignore_x_attrs_error: false
    options:
      tag:
        - documents
        - daily
      exclude:
        - "*.tmp"
        - "*.log"
      skip-if-unchanged: true
    hooks:
      before:
        - echo "Starting backup: $CRESTIC_JOB_NAME"
      success:
        - echo "Backup completed: $CRESTIC_JOB_NAME"
      failure:
        - echo "Backup failed: $CRESTIC_JOB_NAME - $CRESTIC_ERROR" >&2
```

## Running Backup Jobs

### Run Specific Job

```bash
crestic backup --job documents-backup
```

### Run Multiple Jobs

```bash
crestic backup --job documents-backup,photos-backup
```

### Run All Backup Jobs

```bash
crestic backup --all
```

### Dry Run

```bash
crestic backup --job documents-backup --dry-run
```

## Error Handling

When running multiple jobs, each job executes **independently**. If one job fails:

- The error is logged
- Execution continues with the next job
- Other jobs will still run
- At the end, all errors are collected and returned as a combined error message

This ensures that a failure in one job doesn't prevent other jobs from completing.
Each job's success or failure is tracked separately,
and healthcheck notifications are sent for each job individually.

## See Also

- [Copy Job](/jobs/copy) - Copy snapshots between repositories
- [Configuration Guide](/config) - Complete configuration reference
- [Repositories](/repositories) - Repository setup
- [Hooks](/hooks) - Lifecycle hooks
- [Healthchecks](/healthchecks) - Monitoring integration
