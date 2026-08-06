# Scenario

**Feature**: drifted on-disk `SKILL.md` is refreshed on update

```
HandleInstall -> user/tool mutates SKILL.md
HandleUpdate -> overwrite when content differs; polished updated status
```

## Preconditions

- Skill pre-installed with canonical embedded content.
- `PostInstallMutate` simulates drift before update runs.

## Steps

1. Pre-install.
2. Mutate `SKILL.md` on disk.
3. Run update (optionally with `--dry-run` in child leaves).

## Context

- Children differ only by dry-run modifier.
