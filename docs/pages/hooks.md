# 🪝Hooks

Execute custom commands at different stages of backup or copy job.

## Overview

Hooks allow you to run custom scripts or commands before, on success/failure of backup or copy job.

## Available Hooks

- `before` - Run before backup or copy job starts
- `success` - Run after successful backup or copy
- `failure` - Run after failed backup or copy

## Configuration

```yaml
jobs:
  - type: backup
    name: documents
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

- `CRESTIC_JOB_NAME` - Name of the job
- `CRESTIC_EXIT_CODE` - Exit code of the operation
- `CRESTIC_ERROR` - Error message (only in failure hooks)

## Examples

### Database Backup Before Files

```yaml
jobs:
  - type: backup
    name: full-backup
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
jobs:
  - type: backup
    name: important-backup
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
jobs:
  - type: backup
    name: external-drive
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

- [Configuration Guide](/config) - Complete configuration reference
- [Healthchecks](/healthchecks) - Built-in monitoring integration
