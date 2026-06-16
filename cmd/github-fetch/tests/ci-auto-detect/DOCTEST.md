# CI Auto-Detect Tests

Tests for `github-fetch ci` auto-detecting the repository from `git remote get-url origin`
when no URL or PR reference is given. Uses a mock GitHub API server for deterministic results.

## Decision Tree

```
cmd/github-fetch/tests/ci-auto-detect/
├── [git context]
│   ├── not-in-git-repo/              # Not in a git repo → error
│   └── [in git repo]
│       ├── no-origin-remote/         # No origin remote → error
│       └── [origin configured]
│           ├── [args: none]
│           │   └── no-args-auto-detect/    # Lists all runs + detected URL
│           ├── [--workflow <name>]
│           │   ├── with-workflow-no-logs/  # Filtered run list, no logs
│           │   └── no-matching-workflow/   # No matching workflow → error
│           ├── [--logs --workflow <name>]
│           │   └── with-logs-and-workflow/ # Logs for matching workflow
│           ├── [--run-id <id> --logs]
│           │   └── with-run-id/            # Logs for specific run ID
│           └── [--logs only]
│               └── logs-only-no-workflow/  # Logs for latest run
```

## Test Index

| # | Test Leaf | Description |
|---|-----------|-------------|
| 1 | no-args-auto-detect | `github-fetch ci` with no args auto-detects repo from git origin and lists runs |
| 2 | not-in-git-repo | Error when not in a git repo |
| 3 | no-origin-remote | Error when origin remote is missing |
| 4 | with-workflow-no-logs | Lists workflow runs filtered by `--workflow` name |
| 5 | with-logs-and-workflow | Shows logs for matching workflow |
| 6 | with-run-id | Fetches specific run by `--run-id` |
| 7 | no-matching-workflow | Error when workflow name not found |
| 8 | logs-only-no-workflow | Shows logs for latest run without workflow filter |

## How to Run

```sh
doctest test -v ./cmd/github-fetch/tests/ci-auto-detect
```