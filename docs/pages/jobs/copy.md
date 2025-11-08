# 📨 Copy Job

Copy jobs replicate snapshots from one repository to another. This is useful for creating off-site backups or maintaining multiple backup copies.

## Configuration Structure

```yaml
jobs:
  - type: copy
    name: string                    # Required: Unique job name
    from: string                    # Required: Source repository name
    to: string                      # Required: Target repository name
    cron: string                    # Optional: Cron expression
    healthcheck_url: string         # Optional: Job-specific healthcheck URL
    options:                        # Optional: Restic copy options
      key: value
    hooks:                          # Optional: Lifecycle hooks
      before: []string
      success: []string
      failure: []string
```

## Required Fields

### `type`

Must be `"copy"`.

### `name`

Unique identifier for the job. Used in logs and when selecting specific jobs.

```yaml
name: offsite-copy
name: documents-to-remote
name: backup-replication
```

### `from`

Name of the source repository (must be defined in `repositories` section).

```yaml
from: local-repo
from: primary-backup
```

### `to`

Name of the target repository (must be defined in `repositories` section).

```yaml
to: remote-repo
to: secondary-backup
```

## Optional Fields

### `cron`

Cron expression for scheduling automated copy operations.

**Format**: `minute hour day month weekday`

**Examples**:
```yaml
cron: "0 3 * * *"      # Daily at 3:00 AM (after backup completes)
cron: "0 */12 * * *"  # Every 12 hours
cron: "0 4 * * 0"     # Weekly on Sunday at 4:00 AM
```

See [Cron Command](/cli/cron) for more details.

### `healthcheck_url`

Override global healthcheck URL for this specific job.

```yaml
healthcheck_url: https://hc-ping.com/uuid-for-copy/copy-job
```

See [Healthchecks](/healthchecks) for more details.

## Options

The `options` field accepts any restic copy option. Common options:

### Filter by Tags

Copy only snapshots with specific tags:

```yaml
options:
  tag:
    - important
    - documents
    - daily
```

### Filter by Hostname

Copy only snapshots from specific host:

```yaml
options:
  host: my-server
```

### Filter by Paths

Copy only snapshots containing specific paths:

```yaml
options:
  path: /home/user/Documents
```

**See**: [Restic Copy Documentation](https://restic.readthedocs.io/en/stable/045_working_with_repos.html#copying-snapshots-between-repositories) for complete list of options.

## Hooks

Execute custom commands at different stages:

```yaml
hooks:
  before:
    - echo "Starting copy operation: $CRESTIC_JOB_NAME"
  success:
    - echo "Copy completed successfully: $CRESTIC_JOB_NAME"
  failure:
    - echo "Copy failed: $CRESTIC_JOB_NAME - $CRESTIC_ERROR" >&2
```

**Environment variables available in hooks**:
- `CRESTIC_JOB_NAME` - Name of the job
- `CRESTIC_EXIT_CODE` - Exit code of the operation
- `CRESTIC_ERROR` - Error message (only in failure hooks)

See [Hooks](/hooks) for more details.

## Complete Example

```yaml
jobs:
  - type: copy
    name: documents-copy-to-remote
    from: local-repo
    to: remote-repo
    cron: "0 3 * * *"
    healthcheck_url: https://hc-ping.com/uuid/copy-job
    options:
      tag:
        - documents
        - important
      host: my-server
    hooks:
      before:
        - echo "Starting copy: $CRESTIC_JOB_NAME"
      success:
        - echo "Copy completed: $CRESTIC_JOB_NAME"
      failure:
        - echo "Copy failed: $CRESTIC_JOB_NAME - $CRESTIC_ERROR" >&2
```

## Running Copy Jobs

### Run Specific Job

```bash
crestic backup --job offsite-copy
```

**Note**: Copy jobs are executed using the `backup` command, not a separate `copy` command.

### Run Multiple Jobs

```bash
crestic backup --job offsite-copy,another-copy
```

### Run All Jobs

```bash
crestic backup --all
```

## What Happens During Copy

The copy operation:

1. **Sends start ping** to healthcheck service (if configured)
2. **Runs 'before' hooks** (if configured)
3. **Copies snapshots** from source to target repository
4. **Runs 'success' or 'failure' hooks** based on outcome
5. **Sends success/failure ping** to healthcheck service

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

- [Backup Job](/jobs/backup) - Back up directories to repositories
- [Configuration Guide](/config) - Complete configuration reference
- [Repositories](/repositories) - Repository setup
- [Hooks](/hooks) - Lifecycle hooks
- [Healthchecks](/healthchecks) - Monitoring integration
