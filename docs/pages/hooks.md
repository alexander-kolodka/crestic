# 🪝Hooks

Execute custom commands at different stages of backup or copy job.

## Overview

Hooks allow you to run custom scripts or commands before, on success/failure of backup or copy job. Hooks are configured on individual jobs within a pipeline.

## Available Hooks

- `before` - Run before backup or copy job starts
- `success` - Run after successful backup or copy
- `failure` - Run after failed backup or copy

## Configuration

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

Hooks have access to these environment variables:

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
- If a `before` hook fails (non-zero exit code), the job is aborted

## See Also

- [Pipelines](/pipelines) - Pipeline configuration
- [Configuration Guide](/config) - Complete configuration reference
- [Healthchecks](/healthchecks) - Built-in monitoring integration
