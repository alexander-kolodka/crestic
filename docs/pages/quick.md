# 🚀 Quick Start

## Installation

### Install via Homebrew

```bash
brew tap alexander-kolodka/crestic
brew install crestic
```

### Install from source (requires Go 1.25+)

```bash
go install github.com/alexander-kolodka/crestic@latest
```

## Basic Configuration

Create a `crestic.yaml` file:

```yaml
# Global healthcheck URL (optional)
healthcheck_url: https://hc-ping.com/your-uuid-here

jobs:
  - type: backup
    name: documents
    from:
      - /home/user/Documents
      - /home/user/Projects
    to: local-repo
    cron: "0 2 * * *"  # Daily at 2 AM
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
# Run a backup (auto-initializes repository if needed)
crestic backup --job documents

# What happens during backup:
# 1. Checks/initializes repository
# 2. Creates encrypted snapshot
# 3. Verifies integrity (check)
# 4. Applies retention policy (forget)

# Run all scheduled jobs (use this in system cron)
crestic cron
```

## Schedule with System Cron

Add to your crontab:

```cron
# Check for scheduled backups every 5 minutes
*/5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml
```

Crestic keeps track of the last run time,
so even if it’s executed infrequently, it won’t skip any scheduled jobs.
