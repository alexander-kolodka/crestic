# ✅ Check

```bash
crestic check [--all, -a] [--repo, -r <name>]
```

Check and initialize repositories.

## Flags

- `--all, -a` - Check all repositories
- `--repo, -r <name>` - Check specific repository/repositories

## Examples

```bash
# Check all repositories
crestic check --all

# Check specific repository
crestic check --repo local-repo

# Check multiple repositories
crestic check --repo local-repo,remote-repo
```

## Behavior

For each repository:
1. Checks if repository is initialized
2. If not initialized, creates new repository
3. If initialized, verifies repository integrity

**Note**: The `backup` command automatically runs check, so you usually don't need to run this separately.
