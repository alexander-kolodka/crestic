# 📨 Copy Job

Copy jobs replicate snapshots from one repository to another. This is useful for creating off-site backups or maintaining multiple backup copies. Jobs are defined inside a [pipeline](/pipelines).

## Configuration Structure

```yaml
pipelines:
  - name: my-pipeline
    jobs:
      - type: copy
        name: string                    # Required: Job name (unique within pipeline)
        from: string                    # Required: Source repository name
        to: string                      # Required: Target repository name
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

Identifier for the job within the pipeline. Used in logs and CLI as `pipeline/job`.

```yaml
name: offsite-copy
name: documents-to-remote
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
- `CRESTIC_JOB_NAME` - Qualified job name (`pipeline/job`)
- `CRESTIC_EXIT_CODE` - Exit code of the operation
- `CRESTIC_ERROR` - Error message (only in failure hooks)

See [Hooks](/hooks) for more details.

## Complete Example

```yaml
pipelines:
  - name: documents-nightly
    cron: "0 2 * * *"
    jobs:
      - type: backup
        name: local-backup
        from: [/home/user/Documents]
        to: local-repo

      - type: copy
        name: offsite-copy
        from: local-repo
        to: remote-repo
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

### Run Entire Pipeline

```bash
crestic backup --pipeline documents-nightly
```

### Run Specific Copy Job

```bash
crestic backup --job documents-nightly/offsite-copy
```

**Note**: Copy jobs are executed using the `backup` command, not a separate `copy` command.

### Run All Jobs

```bash
crestic backup --all
```

## What Happens During Copy

The copy operation:

1. **Sends start ping** to healthcheck service (if configured)
2. **Runs `before` hooks** (if configured)
3. **Copies snapshots** from source to target repository
4. **Runs `success` or `failure` hooks** based on outcome
5. **Sends success/failure ping** to healthcheck service

## Error Handling

Jobs in a list run **sequentially** (fail-fast). If one job fails:

- The error is logged and returned immediately
- Remaining jobs in the same list are not executed

When running multiple pipelines, each pipeline still runs even if a previous pipeline
failed; pipeline-level errors are combined at the end.

## See Also

- [Pipelines](/pipelines) - Pipeline configuration and scheduling
- [Backup Job](/jobs/backup) - Back up directories to repositories
- [Configuration Guide](/config) - Complete configuration reference
- [Repositories](/repositories) - Repository setup
- [Hooks](/hooks) - Lifecycle hooks
- [Healthchecks](/healthchecks) - Monitoring integration
