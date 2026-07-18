# List PRs Tests

Tests for `github-fetch prs` (and `pr --list` alias): paginated listing of pull
requests with state filtering, repo auto-detection, and pagination footer hints.

## Decision Tree

```
cmd/github-fetch/tests/list-prs/
├── [repo resolution]
│   ├── repo-resolution/
│   │   ├── not-in-git-repo/          # prs without repo, outside git → error
│   │   └── explicit-repo/            # prs owner/repo without git → success
│   └── [git repo + auto-detect]
│       └── git-repo/
│           ├── [command variant]
│           │   └── command/
│           │       ├── open-default/       # prs, open state, page 1
│           │       └── pr-list-alias/      # pr --list, same output
│           ├── [pagination]
│           │   └── pagination/
│           │       ├── page-2/             # --page 2 --per-page 2
│           │       └── pagination-footer/  # Link header → --page 2 hint
│           ├── state-closed/               # --state closed
│           └── empty-results/              # no PRs → success message
```

## Test Index

| # | Test Leaf | Description |
|---|-----------|-------------|
| 1 | open-default | `prs` in git repo lists open PRs page 1 with auto-detected repo |
| 2 | pr-list-alias | `pr --list` produces same output as `open-default` |
| 3 | explicit-repo | `prs owner/repo` works without git repo |
| 4 | page-2 | `prs --page 2 --per-page 2` shows second page items |
| 5 | state-closed | `prs --state closed` lists only closed PRs |
| 6 | pagination-footer | Page 1 footer hints `--page 2` when more results exist |
| 7 | not-in-git-repo | Error when not in git repo and no explicit repo |
| 8 | empty-results | Success with empty message when no PRs match |

## How to Run

```sh
doctest test -v ./cmd/github-fetch/tests/list-prs
```