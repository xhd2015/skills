# Scenario

**Feature**: multi-flag update touches only installed target dirs

```
HandleInstall --codex -> SKILL.md under .codex only
HandleUpdate --codex --opencode -> codex updated/skipped, opencode absent
```

## Preconditions

- Install and update use combined location flags.

## Steps

1. Leaves configure which flags are used at install vs update time.

## Context

- Uninstalled targets for extra flags must not be created during update.
