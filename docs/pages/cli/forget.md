# 🗑️ Forget

```bash
crestic forget [--all, -a] [--repo, -r <name>] [--dry-run] [--prune]
```

Remove old snapshots according to retention policy.

## Flags

- `--all, -a` - Run forget for all repositories
- `--repo, -r <name>` - Run forget for specific repository/repositories
- `--dry-run` - Show what would be deleted without deleting
- `--prune` - Actually remove data from repository (frees space)

## Examples

```bash
# Show what would be deleted (safe)
crestic forget --all --dry-run

# Delete old snapshots (marks for deletion)
crestic forget --all

# Delete old snapshots AND free disk space
crestic forget --all --prune
```

## Retention Policy

Configure automatic snapshot retention with `forget_options` in your repository configuration:

```yaml
repositories:
  my-repo:
    forget_options:
      keep-daily: 7
      keep-weekly: 4
      keep-monthly: 12
      prune: true  # Actually frees disk space
```

**Note**: The `backup` command automatically runs forget after each backup, so you usually don't need to run this separately.
