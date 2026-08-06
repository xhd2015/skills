# Scenario

**Feature**: `--global` update uses `$HOME` skill trees

```
HandleInstall --global -> ~/.<tool>/skills/<name>
HandleUpdate --global -> only global tree updated
```

## Preconditions

- Leaves set `req.UseGlobalHome` so `HOME` is an isolated temp dir.

## Steps

1. Pre-install with `--global`.
2. Update with `--global`.

## Context

- Project-local `.agents/...` under workdir must stay untouched.
