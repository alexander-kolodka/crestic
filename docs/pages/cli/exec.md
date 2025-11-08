# 📟 Exec

```bash
crestic exec [--all, -a] [--repo, -r <name>] <command> [-- native-options]
```

Execute native restic commands on repositories.

## Flags

- `--all, -a` - Execute on all repositories
- `--repo, -r <name>` - Execute on specific repository/repositories

## Examples

```bash
# List snapshots for all repositories
crestic exec --all snapshots

# List snapshots for specific repository
crestic exec --repo local-repo snapshots

# List snapshots with native restic options
crestic exec --all snapshots --compact

# Show specific snapshot contents
crestic exec --repo local-repo ls latest

# Get repository stats
crestic exec --repo local-repo stats

# Mount repository (requires FUSE)
crestic exec --repo local-repo mount /mnt/restic
```
