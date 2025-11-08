# Crestic

Crestic is a simple wrapper around the excellent backup tool [restic](https://restic.net/).

Restic is fast, secure, and incredibly capable —
but once you need to back up several folders to multiple destinations,
its CLI can quickly become repetitive and hard to manage.

Crestic makes this easier. You define all your repositories and backup jobs
in a single configuration file and run them with simple commands.
No hidden logic or magic — just more convenience and fewer chances
to mess things up while using restic.

➡️📝 **[Full documentation](https://crestic.kolodla.fyi)** 

## Native features of restic

-	Incremental deduplicated backups — minimal space usage
-	End-to-end encryption
-	Multiple storage backends — local, SFTP, AWS S3, Backblaze B2, rclone remotes, etc.
-	Snapshot-based system with automatic pruning policies
-	Reliable restores from any snapshot
-	Exclude patterns and ignore rules
-	No central server or database required
-	Cross-platform: macOS, Linux, Windows

## What Crestic adds on top

- Single YAML config for all repositories and jobs
- Built-in [healthchecks.io](https://healthchecks.io) support
- Before / after hooks — run custom commands around backup tasks
- Password-command support — pull credentials from keychain, pass, scripts, etc.
- Cron-ready execution

## Installation

```bash
# Install from source (requires Go 1.25+)
go install alexander-kolodka/crestic@latest
```
