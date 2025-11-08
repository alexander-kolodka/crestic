# 🌡️ Healthchecks

Monitor your backups with [Healthchecks.io](https://healthchecks.io) integration.

## Overview

Crestic integrates with Healthchecks.io to notify you when:
- Backups complete successfully
- Backups fail
- Scheduled backups don't run (dead man's switch)

[Healthchecks.io](https://healthchecks.io) is a cron monitoring service that alerts you when jobs fail or don't run on schedule.

## How It Works

1. Crestic sends "start" ping when job begins
2. Sends "success" ping when job completes successfully
3. Sends "failure" ping if job fails
4. Healthchecks alerts you if expected ping doesn't arrive

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

Apply to all jobs:

```yaml
healthcheck_url: https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef

jobs:
  - type: backup
    name: documents
    # ... rest of config
```

#### Per-Job Configuration

Different healthcheck for each job:

```yaml
jobs:
  - type: backup
    name: documents
    healthcheck_url: https://hc-ping.com/uuid-for-documents
    from: [/home/user/Documents]
    to: local-documents-repo

  - type: backup
    name: photos
    healthcheck_url: https://hc-ping.com/uuid-for-photos
    from: [/home/user/Photos]
    to: local-photos-repo
```

#### With Slug

Add slug to ping URL for better identification:

```yaml
healthcheck_url: https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef/documents-backup
```

## See Also

- [Hooks Guide](/hooks) - Custom notifications with hooks
- [Healthchecks.io Documentation](https://healthchecks.io/docs/) - Official docs
