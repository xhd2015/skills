# CI workflow note

## Status

**updated** (workflow aligned with doctest pattern, committed, and pushed)

## Branch / remote / push

| Field | Value |
|-------|--------|
| Branch | `master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Remote | `ssh://git@github.com/xhd2015/skills` (`origin`) |
| Upstream | `origin/master-2026-08-06-use-go-best-practice-to-review-current-project` (set on push) |
| CI commit SHA | `7750ed6a61074bf88988db58b0715525d7f4e821` (`7750ed6`) — workflow + helper |
| Tip SHA (note) | `8b9ea1afa51e640c568c68b973ff6caccf9cd2ec` (`8b9ea1a`) — this note |
| Push result | **success** — `c2a86cf..7750ed6` then `7750ed6..8b9ea1a` on `master-2026-08-06-use-go-best-practice-to-review-current-project` |

## Paths changed (CI commit)

- `.github/workflows/test.yml` — coverage + doctest discovery/e2e + xgo merge + summary + artifacts; kept skills-specific git/Node/Playwright setup
- `script/ci/coverage-package-table.py` — product package table for `GITHUB_STEP_SUMMARY` (module `github.com/xhd2015/skills/`)

## How this differs from doctest’s workflow

| Aspect | doctest | skills (this push) |
|--------|---------|---------------------|
| COVERPKG | `github.com/xhd2015/doctest/...` | `github.com/xhd2015/skills/...` |
| Install doctest | `go install ./cmd/doctest` | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Extra setup | none | git identity; Node 20; Playwright OS libs |
| Stages | go test, discovery `!e2e`, e2e | same three stages |
| Package table | skip `script/`, `cmd/`, `legacy_*` | skip `cmd/` only under skills module |
| e2e leaves | many labeled e2e | none yet; e2e stage still runs for pattern parity (empty profile skipped on merge) |

## View Actions for this push

- Repo: https://github.com/xhd2015/skills  
- Actions (branch filter): https://github.com/xhd2015/skills/actions?query=branch%3Amaster-2026-08-06-use-go-best-practice-to-review-current-project  
- Workflow file: `.github/workflows/test.yml` (workflow name: **Test**)  
- Commit: https://github.com/xhd2015/skills/commit/7750ed6a61074bf88988db58b0715525d7f4e821  
