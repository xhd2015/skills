---
name: github-fetch
description: >-
  Fetch GitHub PR content and create git worktrees for PR development.
  Use when the user wants to inspect a pull request or set up a local
  worktree to work on a PR.
---

# GitHub Fetch Skill

A CLI tool for fetching GitHub PR content and managing PR worktrees.

## Commands

### pr — Fetch and display PR content

Fetches PR metadata and diff from GitHub and prints a human-readable summary.

```bash
github-fetch pr https://github.com/xhd2015/xgo/pull/379
github-fetch pr 379    # auto-detect owner/repo from git remote
```

### work-on — Create a git worktree for a PR

Verifies the current repo matches the PR's base repository, fetches the PR branch, and creates a git worktree.

```bash
github-fetch work-on https://github.com/xhd2015/xgo/pull/379
github-fetch work-on 379
```

### push — Push to the PR's source branch

Pushes the current HEAD to the PR's head branch. Accepts a URL, number, or auto-detects from the current branch name (must match `pr-<number>`).

```bash
github-fetch push https://github.com/xhd2015/xgo/pull/379
github-fetch push 379
github-fetch push -f    # force push, auto-detect PR from branch name
```

### skill show — Show this SKILL.md

```bash
github-fetch skill show
```

### skill install — Install this skill

```bash
github-fetch skill install              # install to .agents/skills/github-fetch
github-fetch skill install --cursor     # install to .cursor/skills/github-fetch
github-fetch skill install ./some/dir   # install to a custom directory
```
