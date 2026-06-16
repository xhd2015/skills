# List Issues Tests

Tests for `github-fetch issues` (and `issue --list` alias): paginated listing of
issues with PR exclusion, state filtering, and repo auto-detection.

## Decision Tree

```
cmd/github-fetch/tests/list-issues/
├── [repo resolution]
│   ├── repo-resolution/
│   │   ├── not-in-git-repo/          # issues without repo, outside git → error
│   │   └── explicit-repo/            # issues owner/repo without git → success
│   └── [git repo + auto-detect]
│       └── git-repo/
│           ├── [command variant]
│           │   └── command/
│           │       ├── open-default/       # issues, PRs excluded from output
│           │       └── issue-list-alias/   # issue --list, same output
│           ├── [pagination]
│           │   └── pagination/
│           │       └── page-2/             # --page 2 --per-page 2
│           ├── state-closed/               # --state closed
│           └── empty-results/              # no issues → success message
```

## Test Index

| # | Test Leaf | Description |
|---|-----------|-------------|
| 1 | open-default | `issues` in git repo lists open issues, excludes PRs |
| 2 | issue-list-alias | `issue --list` produces same output as `open-default` |
| 3 | explicit-repo | `issues owner/repo` works without git repo |
| 4 | page-2 | `issues --page 2 --per-page 2` shows second page items |
| 5 | state-closed | `issues --state closed` lists only closed issues |
| 6 | not-in-git-repo | Error when not in git repo and no explicit repo |
| 7 | empty-results | Success with empty message when no issues match |

## How to Run

```sh
doctest test -v ./cmd/github-fetch/tests/list-issues
```