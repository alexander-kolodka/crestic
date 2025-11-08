# 💾 Backup

```bash
crestic backup [--all, -a] [--job, -j <name>] [--dry-run]
```

Performs a backup of all jobs if the `-a` or `--all` flag is passed. To only backup some jobs pass one or more `-j` or `--job` flags.

The `--dry-run` flag will do a dry run showing what would have been backed up, but won't touch the actual data.

## What It Does

The `backup` command performs a complete backup workflow (all steps are automatic):

1. **Sends start ping** to healthcheck service (if configured)
2. **Runs 'before' hooks** (if configured)
3. **Checks repository** - automatically initializes if not exists
4. **Creates backup** - encrypted, deduplicated snapshot
5. **Verifies integrity** - runs `restic check` on repository
6. **Applies retention policy** - runs `restic forget` with `forget_options`
7. **Runs 'success' or 'failure' hooks** based on outcome
8. **Sends success/failure ping** to healthcheck service

```bash
# All jobs
crestic backup --all

# Specific job
crestic backup --job documents

# Multiple jobs
crestic backup --job documents,photos

# Dry run
crestic backup --all --dry-run
```

## Automatic Cleanup

If your repository has `forget_options` configured, old snapshots are automatically removed after each backup:

```yaml
repositories:
  my-repo:
    forget_options:
      keep-daily: 7
      keep-weekly: 4
      prune: true  # Actually frees disk space
```

For more options, see [Removing backup snapshots](https://restic.readthedocs.io/en/stable/060_forget.html).

## Error Handling

When running multiple jobs, each job executes **independently**. If one job fails:

- The error is logged
- Execution continues with the next job
- Other jobs will still run
- At the end, all errors are collected and returned as a combined error message

This ensures that a failure in one job doesn't prevent other jobs from completing.
Each job's success or failure is tracked separately,
and healthcheck notifications are sent for each job individually.
