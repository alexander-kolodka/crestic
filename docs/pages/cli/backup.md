# 💾 Backup

```bash
crestic backup [--all, -a] [--pipeline, -p <name>] [--job, -j <pipeline/job>] [--dry-run]
```

Performs backup and copy jobs from your configuration.
`backup` supports two execution modes:

- `--pipeline` / `--all`: runs full pipelines (all jobs in pipeline order)
- `--job`: runs only explicitly selected jobs

## What It Does

The `backup` command performs a complete backup workflow (all steps are automatic):

1. **Sends start ping** to healthcheck service (if configured)
2. **Runs `before` hooks** (if configured)
3. **Checks repository** - automatically initializes if not exists
4. **Creates backup** - encrypted, deduplicated snapshot
5. **Verifies integrity** - runs `restic check` on repository
6. **Applies retention policy** - runs `restic forget` with `forget_options`
7. **Runs `success` or `failure` hooks** based on outcome
8. **Sends success/failure ping** to healthcheck service

## Examples

```bash
# All jobs from all pipelines
crestic backup --all

# Entire pipeline
crestic backup --pipeline documents-nightly

# Specific job (qualified name: pipeline/job)
crestic backup --job documents-nightly/local-backup

# Multiple jobs
crestic backup --job documents-nightly/local-backup,photos-weekly/backup

# Dry run
crestic backup --all --dry-run
```

Only one of `--all`, `--pipeline`, or `--job` can be specified at a time.

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
