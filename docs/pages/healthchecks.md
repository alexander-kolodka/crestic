# 🌡️ Healthchecks

Monitor your backups with [Healthchecks.io](https://healthchecks.io) integration.

## Overview

Crestic integrates with Healthchecks.io to notify you when:
- Pipelines complete successfully
- Pipelines fail
- Scheduled pipelines don't run (dead man's switch)

[Healthchecks.io](https://healthchecks.io) is a cron monitoring service that alerts you when jobs fail or don't run on schedule.

Each pipeline can have its own check. Match the Healthchecks.io schedule to that pipeline's `cron`.

## Setup

### 1. Create a Check per Pipeline

In Healthchecks dashboard:

1. Click "Add Check"
2. Set name (e.g., "Crestic documents pipeline")
3. Configure schedule to match the pipeline cron expression
4. Set grace time (how long to wait before alerting)
5. Copy the ping URL

Example ping URL:
```
https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef
```

### 2. Configure Crestic

Add `healthcheck_url` on the pipeline:

```yaml
pipelines:
  - name: documents-nightly
    cron: "0 2 * * *"
    healthcheck_url: https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef
    jobs:
      - type: backup
        name: local-backup
        # ... rest of config
```

#### With Slug

```yaml
pipelines:
  - name: documents-nightly
    cron: "0 2 * * *"
    healthcheck_url: https://hc-ping.com/01234567-89ab-cdef-0123-456789abcdef/documents-nightly
    jobs:
      # ...
```

### 3. Enable Healthchecks

Pings are off by default. Pass `--healthcheck` to enable them:

```bash
# Healthcheck disabled by default
crestic backup --pipeline documents-nightly

# Enable healthcheck explicitly
crestic backup --pipeline documents-nightly --healthcheck
crestic backup --all --healthcheck
```

```bash
# Cron with healthchecks (typical crontab entry)
*/5 * * * * crestic cron --config /path/to/crestic.yaml --healthcheck
```

Notes:
- `--job` runs do not send healthcheck pings (partial pipeline)
- Dry-run does not send pings
- Healthcheck ping failures are logged and do not fail the pipeline

## Ping Payload

Bodies are plain text (not JSON).

**Start** (`/start`):
```
pipeline: documents-nightly
jobs:
  - local-backup
  - offsite-copy
```

**Success**: empty body

**Fail** (`/fail`): error message text

## See Also

- [Hooks Guide](/hooks) - Custom notifications with hooks
- [Pipelines](/pipelines) - Pipeline configuration
- [Healthchecks.io Documentation](https://healthchecks.io/docs/) - Official docs
