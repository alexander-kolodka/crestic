# 🔄 Restore

```bash
crestic restore --repo <name> --target <path> [--snapshot <id>]
```

Restore a snapshot to a directory.

## Flags

- `--repo, -r <name>` - Required: Repository to restore from
- `--target, -t <path>` - Required: Directory to restore to
- `--snapshot, -s <id>` - Snapshot ID (default: "latest")

## Examples

```bash
# Restore latest snapshot
crestic restore --repo local-repo --target ./restore

# Restore specific snapshot
crestic restore --repo local-repo --target ./restore --snapshot abc123
```

## Finding Snapshots

First, list available snapshots:

```bash
crestic exec --repo local-repo snapshots
```

Output shows snapshot IDs:
```
ID        Time                 Host        Tags        Paths
----------------------------------------------------------------
abc123    2024-01-15 10:00:00  myhost      daily       /home/user
def456    2024-01-14 10:00:00  myhost      daily       /home/user
```

Then restore specific snapshot:

```bash
crestic restore --repo local-repo --snapshot abc123 --target ./restore
```
