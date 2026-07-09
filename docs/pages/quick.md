# 🚀 Quick Start

## Installation

### Install via Homebrew

```bash
brew tap alexander-kolodka/crestic
brew install crestic
```

### Install from source (requires Go 1.26+)

```bash
go install github.com/alexander-kolodka/crestic@latest
```

## Basic Configuration

Create a `crestic.yaml` file:

```yaml
# Healthcheck URL (optional)
healthcheck_url: https://hc-ping.com/your-uuid-here

pipelines:
  - name: documents-nightly
    cron: "0 2 * * *"  # Daily at 2 AM
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
        hooks:
          before:
            - echo "Starting backup..."
          success:
            - echo "Backup completed!"
          failure:
            - echo "Backup failed!" >&2

repositories:
  local-repo:
    path: /backup/restic/documents
    password_command: "security find-generic-password -a restic-password -s crestic -w"
    forget_options:
      keep-daily: 7
      keep-weekly: 4
      keep-monthly: 12
```

## Run Your First Backup

```bash
# Run a specific job (qualified name: pipeline/job)
crestic backup --job documents-nightly/local-backup

# Or run the entire pipeline
crestic backup --pipeline documents-nightly

# What happens during backup:
# 1. Checks/initializes repository
# 2. Creates encrypted snapshot
# 3. Verifies integrity (check)
# 4. Applies retention policy (forget)

# Run all scheduled pipelines (use this in system cron)
crestic cron
```

## Schedule with System Cron

Add to your crontab:

```cron
# Check for scheduled pipelines every 5 minutes
*/5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml
```

Crestic keeps track of per-pipeline last run times,
so even if it's executed infrequently, it won't skip any scheduled pipelines.
