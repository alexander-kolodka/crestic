# Crestic

Go CLI wrapper around [restic](https://restic.net/). Users define repositories and backup pipelines in a single YAML config file.

- Go 1.26+
- cobra (CLI), zerolog/slog (logging), testify (tests)
- Docs: https://crestic.kolodka.fyi

## Architecture

```
cmd/ → internal/cases/{feature}/ → entity | dto | restic | shell | healthchecks | hooks
```

- **cmd/** — cobra commands, flag parsing, wiring handlers
- **internal/cases/** — use case handlers (`Handler` + `Command` per feature)
- **internal/cases/handler/** — middleware chain (lock, panic recovery)
- **internal/entity/** — domain types (Job, Pipeline, Config, Repository)
- **internal/dto/** — YAML config parsing and mapping to entity
- **internal/restic/**, **shell/**, **healthchecks/**, **hooks/** — infrastructure

Config structure: repositories + pipelines. Each pipeline contains jobs (backup, copy) with optional cron schedule and hooks.

## Before finishing any task

1. Run `golangci-lint run` — must pass (includes `noinlineerr`, strict revive, sloglint)
2. Run unit tests: `go test $(go list ./... | grep -v '/tests$')`
3. If CLI or config changed: update `docs/pages/` and `cmd/config.example.yaml`
4. If new CLI command added: add integration test in `tests/`

## Reference files

Use these as patterns when generating code:

| Pattern | File |
|---------|------|
| Error handling (split style) | `internal/dto/mapper.go` |
| Handler chain wiring | `cmd/cron.go` |
| Batch pipeline execution | `internal/cases/runpipelines/handler.go` |
| Integration test harness | `tests/harness/sandbox.go` |
| Integration test example | `tests/backup_test.go` |

## Go style

See `.cursor/rules/go-style.mdc` and `.golangci.yml`. Key rules:

- No inline error checks: assign `err` on one line, check on the next
- Package names: single word, no underscores (`runpipelines`)
- Static errors: `errors.New()`, not `fmt.Errorf()` without formatting
- nolint directives must include an explanation

## Commits

- Imperative mood, no trailing period in subject
- Run linter before committing to avoid separate "Fix linter issues" commits
