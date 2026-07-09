# 🔗 Pipelines

Pipelines group related backup and copy jobs into a single workflow that runs on a schedule.

## Why Pipelines?

A typical backup scenario involves multiple steps: back up local files, then copy snapshots to remote storage. These steps are logically related and should run together — not as independent scheduled tasks with manually offset cron times.

**Pipeline** — the scenario you schedule and run (e.g. "nightly documents backup: local → remote").

**Job** — one atomic operation inside a pipeline (a backup or copy step).

When a pipeline's cron triggers, all its jobs run **sequentially** in config order.

## Configuration Structure

```yaml
pipelines:
  - name: string          # Required: Unique pipeline name
    cron: string          # Optional: Cron expression for scheduling
    jobs:                 # Required: List of jobs to run
      - type: backup
        name: string
        # ... backup fields
      - type: copy
        name: string
        # ... copy fields
```

## Required Fields

### `name`

Unique identifier for the pipeline. Used in CLI and logs.

```yaml
name: documents-nightly
name: photos-weekly
```

### `jobs`

List of jobs to execute when the pipeline runs. Job names must be unique within a pipeline.

See [Backup Job](/jobs/backup) and [Copy Job](/jobs/copy) for job configuration details.

## Optional Fields

### `cron`

Cron expression for scheduling the entire pipeline.

**Format**: `minute hour day month weekday`

**Examples**:
```yaml
cron: "0 2 * * *"      # Daily at 2:00 AM
cron: "0 */6 * * *"    # Every 6 hours
cron: "30 3 * * 0"     # Weekly on Sunday at 3:30 AM
```

**To use scheduling**:
1. Set cron expression on the pipeline
2. Add to system crontab: `*/5 * * * * crestic cron --config /path/to/crestic.yaml`
3. Crestic tracks per-pipeline state, so system cron can run every 5-30 minutes

See [Cron Command](/cli/cron) for more details.

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
        to: local-repo
        options:
          tag: [documents]

      - type: copy
        name: offsite-copy
        from: local-repo
        to: remote-repo
        options:
          tag: [documents]

repositories:
  local-repo:
    path: /backup/restic
    password_command: "pass show restic/local"
  remote-repo:
    path: rclone:backblaze:backup/restic
    password_command: "pass show restic/remote"
```

## Running Pipelines

### Run Entire Pipeline

```bash
crestic backup --pipeline documents-nightly
```

### Run Specific Job

Jobs are referenced by qualified name: `pipeline/job`

```bash
crestic backup --job documents-nightly/local-backup
crestic backup --job documents-nightly/offsite-copy
```

### Run Multiple Jobs

```bash
crestic backup --job documents-nightly/local-backup,photos-weekly/backup
```

### Run All Jobs from All Pipelines

```bash
crestic backup --all
```

## Migration from Flat Jobs

The top-level `jobs:` key has been replaced by `pipelines:` with nested jobs.

**Before:**
```yaml
jobs:
  - type: backup
    name: documents-backup
    cron: "0 2 * * *"
    from: [/home/user/Documents]
    to: local-repo

  - type: copy
    name: documents-copy-to-remote
    cron: "0 3 * * *"
    from: local-repo
    to: remote-repo
```

**After:**
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
```

Key changes:
- `cron` moves from individual jobs to the pipeline
- Jobs are nested under `pipelines`
- CLI uses `--pipeline` or `--job pipeline/job` instead of `--job job-name`
- `CRESTIC_JOB_NAME` in hooks is now the qualified name (`documents-nightly/local-backup`)

## See Also

- [Backup Job](/jobs/backup) - Backup job configuration
- [Copy Job](/jobs/copy) - Copy job configuration
- [Cron Command](/cli/cron) - Scheduled execution
- [Configuration Guide](/config) - Complete configuration reference
