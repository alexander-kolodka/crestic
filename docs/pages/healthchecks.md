# 🌡️ Healthchecks

Monitor your backups with [Healthchecks.io](https://healthchecks.io) integration.

## Overview

Crestic integrates with Healthchecks.io to notify you when:
- Backups complete successfully
- Backups fail
- Scheduled backups don't run (dead man's switch)

[Healthchecks.io](https://healthchecks.io) is a cron monitoring service that alerts you when jobs fail or don't run on schedule.

## Setup

### 1. Create Check

In Healthchecks dashboard:

1. Click "Add Check"
2. Set name (e.g., "Crestic Daily Backup")
3. Configure schedule to match your cron expression
4. Set grace time (how long to wait before alerting)
5. Copy the ping URL

Example ping URL:
```
https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef
```

### 2. Configure Crestic

Add healthcheck URL to your configuration:

#### Global Configuration

```yaml
healthcheck_url: https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef

jobs:
  - type: backup
    name: documents
    # ... rest of config
```

#### With Slug

Add slug to ping URL for better identification:

```yaml
healthcheck_url: https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef/daily-backups
```

### 3. Enable Healthchecks

Use the `--healthcheck` flag to enable notifications:

#### Manual Backups

```bash
# Healthcheck disabled by default
crestic backup --all

# Enable healthcheck explicitly
crestic backup --all --healthcheck
```

#### Cron Scheduler

```bash
# Healthcheck disabled by default
crestic cron

# Enable healthcheck explicitly
crestic cron --healthcheck
```
## See Also

- [Hooks Guide](/hooks) - Custom notifications with hooks
- [Healthchecks.io Documentation](https://healthchecks.io/docs/) - Official docs
