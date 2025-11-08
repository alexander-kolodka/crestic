# 🎛 Configuration

## Path

By default crestic searches for a `crestic.yaml` file in the current directory, your home folder and your config folder:

- `./crestic.yaml`
- `~/crestic.yaml`
- `~/.crestic/crestic.yaml`
- `~/.config/crestic/crestic.yaml`

You can also specify a custom file with the `-c path/to/some/config.yaml`

## Example configuration

```yaml | crestic.yaml
healthcheck_url: https://hc-ping.com/your-uuid-here

jobs:
  - type: backup
    name: documents
    from:
      - /home/user/Documents
      - /home/user/Projects
    to: local-repo
    cron: "0 2 * * *"
    options:
      tag:
        - documents
        - daily
      exclude:
        - "*.tmp"
        - "*.log"
    hooks:
      before:
        - echo "Starting backup..."
      success:
        - echo "Backup completed!"

  - type: copy
    name: offsite-copy
    from: local-repo
    to: remote-repo
    cron: "0 3 * * *"
    options:
      tag:
        - important

repositories:
  local-repo:
    path: /backup/restic/documents
    password_command: "security find-generic-password -a restic-password -s crestic -w"
    forget_options:
      keep-daily: 7
      keep-weekly: 4
      keep-monthly: 12

  remote-repo:
    path: rclone:backblaze:backup-bucket/restic
    password_command: "pass show restic/remote"
    forget_options:
      keep-last: 5
      keep-daily: 14
```

## Configuration Structure

### Global Settings

- `healthcheck_url` - Global healthcheck URL used for all jobs unless overridden per-job

### Jobs

Jobs define backup and copy operations. See Jobs section for detailed documentation:

- **[Backup Job](/jobs/backup)** - Backs up local directories to a repository
- **[Copy Job](/jobs/copy)** - Copies snapshots between repositories

### Repositories

Repositories define storage locations for backups. Supports all restic backends:
- Local filesystem
- SFTP
- S3 and S3-compatible storage
- Rclone
- Azure Blob Storage
- Google Cloud Storage
- And more!

See [Repositories](/repositories) for detailed configuration.
