---
name: crestic-release
description: >-
  Create Crestic project releases: compare commits since last tag on main, run
  pre-release checks, generate changelog, bump v0.X.Y version, and publish
  GitHub Release. Use when the user asks for a crestic release, Crestic tag,
  minor/major bump, changelog, or whether a Crestic release is needed.
disable-model-invocation: true
---

# Crestic Release

Publish a GitHub Release with a new `v0.X.Y` tag. Homebrew tap updates automatically via `.github/workflows/brew.yaml` on the `release: released` event — **not** on a bare git tag push.

## Version scheme

Format: `v0.MAJOR.MINOR` — first segment is always `0`.

| Bump | Rule | Example |
|------|------|---------|
| Minor | `0.X.Y` → `0.X.(Y+1)` | `v0.4.0` → `v0.4.1` |
| Major | `0.X.Y` → `0.(X+1).0` | `v0.4.0` → `v0.5.0` |

User must specify `minor` or `major` when requesting a release.

## Workflow

```
Task Progress:
- [ ] Step 1: Fetch and verify clean state
- [ ] Step 2: Check for new commits since last tag
- [ ] Step 3: Pre-release checks (lint, unit, integration tests)
- [ ] Step 4: Compute next version
- [ ] Step 5: Generate changelog and short title (English)
- [ ] Step 6: Show summary and get user confirmation
- [ ] Step 7: Create GitHub Release
```

### Step 1 — Prepare

```bash
git fetch origin
git status   # working tree must be clean
```

Base comparison point — **latest tag on `origin/main`**:

```bash
LAST_TAG=$(git describe --tags --abbrev=0 origin/main)
```

### Step 2 — Check if release is needed

```bash
git log "${LAST_TAG}..origin/main" --oneline
```

If output is empty — **stop** with: "No new commits since ${LAST_TAG}; release not needed."

### Step 3 — Pre-release checks (required)

Match [AGENTS.md](../../../AGENTS.md) and [`.github/workflows/ci.yaml`](../../../.github/workflows/ci.yaml):

```bash
# 1. Lint
golangci-lint run

# 2. Unit tests
go test $(go list ./... | grep -v '/tests$')

# 3. Integration tests (require restic in PATH and a built binary)
RESTIC_VERSION=0.19.0
# macOS Apple Silicon: restic_${RESTIC_VERSION}_darwin_arm64.bz2
# macOS Intel:         restic_${RESTIC_VERSION}_darwin_amd64.bz2
# Linux:               restic_${RESTIC_VERSION}_linux_amd64.bz2
curl -fsSL "https://github.com/restic/restic/releases/download/v${RESTIC_VERSION}/restic_${RESTIC_VERSION}_darwin_arm64.bz2" \
  | bunzip2 > /tmp/restic && chmod +x /tmp/restic
export PATH="/tmp:$PATH"

go build -o bin/crestic .
go test -parallel 4 ./tests/...
```

If `restic` is already in PATH, skip the curl step. `go build` and `go test ./tests/...` are always required.

If any step fails — **do not create a release**.

### Step 4 — Compute next version

Parse `LAST_TAG` (e.g. `v0.4.0`):

```bash
VERSION="${LAST_TAG#v}"          # 0.4.0
IFS='.' read -r _ MAJOR MINOR <<< "$VERSION"
# minor: NEW_TAG="v0.${MAJOR}.$((MINOR+1))"
# major: NEW_TAG="v0.$((MAJOR+1)).0"
```

Verify tag does not exist yet:

```bash
git rev-parse "$NEW_TAG"   # must fail
```

### Step 5 — Changelog and short title (English)

Source: `git log "${LAST_TAG}..origin/main"`.

Rules:
1. **Merge PR commits** — use the title from the commit body (line after the blank line), e.g. `Add hooks to pipelines` from PR #18
2. **Regular commits** — use subject as-is (`Add Healthchecks for Pipelines`, `Fix Homebrew tap bump...`)
3. **Skip**: `cursor:`, bare `Merge pull request` lines without useful body, internal refactors with no user-facing effect
4. **Group** by category:

```markdown
## What's Changed

### Features
- Add hooks to pipelines (#18)
- Add Healthchecks for Pipelines (#17)

### Fixes
- Fix Homebrew tap bump failing on HTTP 303 redirect (#13)

### Other
- ...
```

**Short title (`SHORT_TITLE`)**: English, max 4–5 words, derived from the changelog. Prefer the dominant user-facing theme. Examples: `Output format changes`, `Pipeline hooks support`, `Fix Homebrew issues`, `Healthchecks and hooks`. If themes are equal weight, pick the most user-facing one; do not invent marketing fluff.

- **Tag**: `v0.X.Y` (unchanged)
- **GitHub Release title**: `v0.X.Y — ${SHORT_TITLE}` (em dash)

Write notes to a temp file before `gh release create`.

### Step 6 — User confirmation

Show before publishing:
- `LAST_TAG` → `NEW_TAG`
- bump type (minor/major)
- proposed release title (`${NEW_TAG} — ${SHORT_TITLE}`)
- full changelog
- commit SHA on `origin/main`

**Do not publish without explicit user approval.**

### Step 7 — Create release

```bash
gh release create "${NEW_TAG}" \
  --target origin/main \
  --title "${NEW_TAG} — ${SHORT_TITLE}" \
  --notes-file /tmp/release-notes.md
```

`gh release create` creates both the tag and the GitHub Release — this triggers the Homebrew bump workflow.

After creation, report the release URL and that the Homebrew tap will update automatically.

## Safety rules

- Release only from `origin/main`, never from feature branches
- Do not edit `internal/version/version.go` — binary version comes from `-ldflags` at build time
- Do not create `CHANGELOG.md` in the repo — release notes live in GitHub Releases only
- Do not push a tag separately from `gh release create`
- Requires `gh` authenticated with access to `alexander-kolodka/crestic`

## Example

User: "Make crestic minor release"

Agent: check commits → lint + unit + integration tests → `v0.4.0` → `v0.4.1` → changelog + short title `Pipeline hooks support` → confirm → `gh release create` with title `v0.4.1 — Pipeline hooks support`
