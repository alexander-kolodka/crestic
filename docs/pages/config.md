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
pipelines:
  - name: documents-nightly
    cron: "0 2 * * *"
    healthcheck_url: https://hc-ping.com/your-uuid-here
    jobs:
      - type: backup
        name: local-backup
        from:
          - /home/user/Documents
          - /home/user/Projects
        to: local-repo
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

### Pipelines

Pipelines group related jobs and define the schedule. Optional `healthcheck_url` enables [Healthchecks.io](/healthchecks) monitoring per pipeline. See [Pipelines](/pipelines) for detailed documentation.

### Jobs

Jobs define individual backup and copy operations within a pipeline:

- **[Backup Job](/jobs/backup)** - Backs up local directories to a repository
- **[Copy Job](/jobs/copy)** - Copies snapshots between repositories

### Repositories

Repository definitions. See [Repositories](/repositories) for detailed documentation.
