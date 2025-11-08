# 🗂 Repositories

Repositories define storage locations for backups. Crestic supports all restic backends.

## Configuration

```yaml
repositories:
  repository-name:
    path: string                  # Required: Repository path or URL
    password_command: string      # Required: Command to get password
    forget_options:               # Optional: Retention policy
      key: value
```

## Password Management

The `password_command` field specifies a shell command that outputs the repository password. Examples:

### macOS Keychain

```yaml
password_command: "security find-generic-password -a restic-password -s crestic -w"
```

### Linux Secret Service

```yaml
password_command: "secret-tool lookup service restic password crestic"
```

### Pass (Password Store)

```yaml
password_command: "pass show restic/repo-name"
```

### GPG Encrypted File

```yaml
password_command: "gpg --decrypt password.gpg 2>/dev/null"
```

### Environment Variable

```yaml
password_command: "echo \"$RESTIC_PASSWORD\""
```

## Retention Policy

Configure automatic snapshot retention with `forget_options`.
The backup command runs 'restic forget' with these options after every backup.
For more options, see [Removing backup snapshots](https://restic.readthedocs.io/en/stable/060_forget.html).

```yaml
repositories:
  my-repo:
    forget_options:
      keep-last: 10
      keep-daily: 7
      keep-weekly: 4
      keep-monthly: 12
      keep-yearly: 3
      keep-within: "2y5m7d"
      keep-tag: ["important"]
      prune: true  # Actually frees disk space
```

These options are automatically applied after each backup.

## Supported Backends

- Local Filesystem
- REST Server
- S3 and S3-compatible
- Backblaze B2
- Azure Blob Storage
- Google Cloud Storage
- Rclone
- and more...

See [restic documentation](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html) for complete list of supported backends.
