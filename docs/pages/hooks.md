# 🪝Hooks

Execute custom commands at different stages of a pipeline or job.

## Overview

Hooks allow you to run custom scripts or commands before, on success, or on failure.
They can be configured on a **pipeline** (around the whole pipeline run) and/or on
individual **jobs** (around a single backup or copy).

## Available Hooks

- `before` - Run before the pipeline or job starts
- `success` - Run after successful completion
- `failure` - Run when the pipeline or job does not reach success (including a failed `before` hook, and a failed job/pipeline body)

## Pipeline Hooks

On a full pipeline run (`backup --all`, `backup --pipeline`, or `cron`), order is:

1. Healthcheck start ping (if enabled)
2. Pipeline `before` hooks
3. All jobs (each with its own job hooks)
4. Pipeline `success` or `failure` hooks
5. Healthcheck success/failure ping (if enabled)

So healthchecks wrap the pipeline; pipeline hooks wrap the jobs.

```yaml
pipelines:
  - name: documents-nightly
    hooks:
      before:
        - echo "Starting pipeline $CRESTIC_PIPELINE_NAME..."
      success:
        - echo "Pipeline completed"
      failure:
        - echo "Pipeline failed: $CRESTIC_ERROR" >&2
    jobs:
      - type: backup
        name: local-backup
        from: [/home/user/Documents]
        to: local-repo
```

Pipeline hooks run only when the pipeline is executed as a whole (`backup --all`,
`backup --pipeline`, or `cron`). They are **not** run for `backup --job`.

## Job Hooks

Job hooks wrap a single backup or copy job.

```yaml
pipelines:
  - name: documents-nightly
    jobs:
      - type: backup
        name: local-backup
        from: [/home/user/Documents]
        to: local-repo
        hooks:
          before:
            - echo "Starting backup..."
            - /usr/local/bin/snapshot-database.sh
          success:
            - echo "Backup successful!"
            - curl -X POST https://your-webhook.com/success
          failure:
            - echo "Backup failed!" >&2
            - /usr/local/bin/alert-admin.sh
```

## Environment Variables

### Pipeline hooks

- `CRESTIC_PIPELINE_NAME` - Pipeline name (e.g. `documents-nightly`)
- `CRESTIC_ERROR` - Error message (only in failure hooks)

### Job hooks

- `CRESTIC_JOB_NAME` - Qualified job name (`pipeline/job`, e.g. `documents-nightly/local-backup`)
- `CRESTIC_EXIT_CODE` - Exit code of the operation
- `CRESTIC_ERROR` - Error message (only in failure hooks)

## Examples

### Database Backup Before Files

```yaml
pipelines:
  - name: full-backup
    jobs:
      - type: backup
        name: backup
        from: [/home/user]
        to: local-repo
        hooks:
          before:
            - pg_dump mydb > /tmp/mydb.sql
            - mysqldump mydb > /tmp/mydb.sql
          success:
            - rm /tmp/mydb.sql
```

### Custom Notifications

```yaml
pipelines:
  - name: important-backup
    jobs:
      - type: backup
        name: backup
        from: [/important/data]
        to: remote-repo
        hooks:
          success:
            - curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
                -d '{"text":"Backup completed successfully"}'
          failure:
            - curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
                -d '{"text":"Backup failed: $CRESTIC_ERROR"}'
```

### Mount/Unmount Volumes

```yaml
pipelines:
  - name: external-drive
    jobs:
      - type: backup
        name: backup
        from: [/mnt/external]
        to: local-repo
        hooks:
          before:
            - mount /dev/sdb1 /mnt/external
          success:
            - umount /mnt/external
          failure:
            - umount /mnt/external
```

## Exit Codes

- Hooks run sequentially
- If a `before` hook fails (non-zero exit code), the pipeline or job is aborted and `failure` hooks still run
- If a `failure` hook itself fails, the error is logged as a warning and does not change the original failure or exit code
- If a `success` hook fails, the pipeline or job fails (no `failure` hooks are run for that)

## See Also

- [Pipelines](/pipelines) - Pipeline configuration
- [Configuration Guide](/config) - Complete configuration reference
- [Healthchecks](/healthchecks) - Built-in monitoring integration
