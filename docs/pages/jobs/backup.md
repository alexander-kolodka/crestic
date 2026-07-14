# 💾 Backup Job

Backup jobs back up local directories to a restic repository. Jobs are defined inside a [pipeline](/pipelines).

## Configuration Structure

```yaml
pipelines:
  - name: my-pipeline
    jobs:
      - type: backup
        name: string                    # Required: Job name (unique within pipeline)
        from: []string                  # Required: Source directories
        to: string                      # Required: Target repository name
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

Identifier for the job within the pipeline. Used in logs and CLI as `pipeline/job`.

```yaml
name: local-backup
name: photos-daily
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

### `ignore_x_attrs_error`
Some filesystems (e.g. Cryptomator, other FUSE mounts) do not allow reading extended file attributes (xattrs).
When restic encounters such files, it exits with status code 3, which means:

>"incomplete metadata for ${file}"

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
        from:
          - /home/user/Documents
          - /home/user/Projects
        to: local-repo
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

### Run Entire Pipeline

```bash
crestic backup --pipeline documents-nightly
```

### Run Specific Job

```bash
crestic backup --job documents-nightly/local-backup
```

### Run Multiple Jobs

```bash
crestic backup --job documents-nightly/local-backup,photos-weekly/backup
```

### Run All Jobs

```bash
crestic backup --all
```

### Dry Run

```bash
crestic backup --job documents-nightly/local-backup --dry-run
```

## Error Handling

Jobs in a list run **sequentially** (fail-fast). If one job fails:

- The error is logged and returned immediately
- Remaining jobs in the same list are not executed

When running multiple pipelines, each pipeline still runs even if a previous pipeline
failed; pipeline-level errors are combined at the end.

## See Also

- [Pipelines](/pipelines) - Pipeline configuration and scheduling
- [Copy Job](/jobs/copy) - Copy snapshots between repositories
- [Configuration Guide](/config) - Complete configuration reference
- [Repositories](/repositories) - Repository setup
- [Hooks](/hooks) - Lifecycle hooks
- [Healthchecks](/healthchecks) - Monitoring integration
