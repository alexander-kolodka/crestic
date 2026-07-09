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
- Only one instance of `crestic cron` can run per configuration file name
- A lock file is created in `~/.crestic/` and uses only the filename of the config, not the full path or extension
  - Example: `/etc/backup/crestic.yaml` and `/etc/backup/crestic.yml` will share the same lock
  - `my.yml` and `config.yml` will run in parallel, as their filenames differ
- This prevents two processes from running the same pipelines simultaneously, while still allowing multiple independent configs to run at the same time

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

Crestic stores per-pipeline last run times in `~/.crestic/crestic-cron-state.json`:

```json
{
  "pipelines": {
    "documents-nightly": { "last_run": "2026-07-07T02:00:00Z" },
    "photos-weekly": { "last_run": "2026-07-06T03:00:00Z" }
  }
}
```

## Notes

- Crestic stores per-pipeline execution timestamps to ensure pipelines run even if cron wasn't triggered exactly on time (e.g. machine was off)
- If no pipeline is due — it exits without doing anything
- If a job fails, Crestic proceeds to the next job in the pipeline (jobs are independent)
