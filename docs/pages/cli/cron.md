# ⏱️ Cron

```bash
crestic cron
```

Run scheduled pipelines based on cron expressions.

## Description

`crestic cron` is designed to be executed by system schedulers (cron, systemd timers, launchd, etc.).

When launched, it:

- Determines which pipelines should have run since their last execution time
- Passes due pipelines to the pipeline runner use case
- Executes all jobs in each due pipeline sequentially
- Remembers last run time per pipeline (won't miss scheduled runs)

## Locking behavior
- Only one instance of `crestic cron` can run per configuration file
- A lock file is created in `~/.crestic/` with the config basename and an MD5 hash of the config's canonical absolute path
  - Format: `crestic-cron-{basename}-{hash}.lock`
  - The hash is computed from the resolved absolute path after symlink resolution, not from the raw `--config` flag value
  - Example: `~/crestic.yaml`, `./crestic.yaml`, and `/home/user/crestic.yaml` share the same lock if they point to the same file
  - Example: `/etc/prod/crestic.yaml` and `/etc/staging/crestic.yaml` use different lock files even if both basenames are `crestic`
- This prevents two processes from running the same config simultaneously, while still allowing multiple independent configs to run at the same time

## Examples

```bash
# Run scheduled pipelines (typically called from system cron)
crestic cron

# Add to system crontab
*/5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml
```

## Scheduling

Define cron expressions on pipelines, not individual jobs:

```yaml
pipelines:
  - name: documents-nightly
    cron: "0 2 * * *"      # Daily at 2 AM
    jobs:
      - type: backup
        name: local-backup
        # ... rest of config

  - name: photos-weekly
    cron: "0 3 * * 0"      # Weekly on Sunday at 3 AM
    jobs:
      - type: backup
        name: backup
        # ... rest of config
```

## State File

Crestic stores per-pipeline last run times in `~/.crestic/` using one state file per configuration:

- Format: `crestic-cron-state-{basename}-{hash}.json`
- `{basename}` is the config filename without extension
- `{hash}` is an MD5 hash of the config's canonical absolute path (same scheme as the lock file)

Example path: `~/.crestic/crestic-cron-state-crestic-a1b2c3....json`

```json
{
  "pipelines": {
    "documents-nightly": { "last_run": "2026-07-07T02:00:00Z" },
    "photos-weekly": { "last_run": "2026-07-06T03:00:00Z" }
  }
}
```

Pipelines with the same name in different config files keep separate state because each config has its own state file.

### First run

When a pipeline has no saved state yet:

- Crestic runs it only if the most recent scheduled slot was within the last 5 minutes
- Otherwise it records the current time and waits for the next scheduled slot

## Notes

- Crestic stores per-pipeline execution timestamps to ensure pipelines run even if cron wasn't triggered exactly on time (e.g. machine was off)
- If no pipeline is due — it exits without doing anything
- If a job fails, remaining jobs in that pipeline are not executed (fail-fast); other due pipelines still run
