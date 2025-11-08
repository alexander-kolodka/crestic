# 🔓 Unlock

```bash
crestic unlock [--all, -a] [--repo, -r <name>]
```

Remove stale locks from repositories.

## Flags

- `--all, -a` - Unlock all repositories
- `--repo, -r <name>` - Unlock specific repository/repositories

## Examples

```bash
# Unlock all repositories
crestic unlock --all

# Unlock specific repository
crestic unlock --repo local-repo
```

## When to Use

Repositories are automatically locked during operations. Sometimes locks are not released if:
- Process is killed
- System crashes
- Network interruption

This command forcefully removes stale locks.

⚠️ **Warning**: Never unlock while another operation is running, or data corruption may occur.
