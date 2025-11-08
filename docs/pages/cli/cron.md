# ⏱️ Cron

```bash
crestic cron
```

Run scheduled jobs based on cron expressions.

## Description

`crestic cron` is designed to be executed by system schedulers (cron, systemd timers, launchd, etc.).

When launched, it:

- Determines which jobs should have run since the last execution time
- Executes all due jobs in the correct order
- Remembers last run time (won't miss jobs)

## Locking behavior
- Only one instance of `crestic cron` can run per configuration file name
- A lock file is created in `~/.crestic/` and uses only the filename of the config, not the full path or extension
  - Example: `/etc/backup/crestic.yaml` and `/etc/backup/crestic.yml` will share the same lock
  - `my.yml` and `config.yml` will run in parallel, as their filenames differ
- This prevents two processes from running the same jobs simultaneously, while still allowing multiple independent configs to run at the same time

## Examples

```bash
# Run scheduled jobs (typically called from system cron)
crestic cron

# Add to system crontab
*/5 * * * * /usr/local/bin/crestic cron --config /path/to/crestic.yaml
```

## Scheduling

Define cron expressions in job configuration:

```yaml
jobs:
  - name: daily-backup
    cron: "0 2 * * *"      # Daily at 2 AM
    # ... rest of config

  - name: hourly-backup
    cron: "0 * * * *"      # Every hour
    # ... rest of config
```

Notes

- Crestic stores the last execution timestamp to ensure jobs are run even if cron wasn’t triggered exactly on time (e.g. machine was off)
- If no job is due — it exits without doing anything
- If a job fails, Crestic proceeds to the next one (jobs are independent)
