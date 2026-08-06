# Scenario

**Feature**: registry batch reports polished status for every skill

```
HandleUpdateMany -> for each skill in registry order:
  installed? -> up to date | updated | would update (+ file lines)
  else -> name  not installed
-> blank line + summary
```

## Preconditions

- `ManySkills` defaults to `skill-alpha` and `skill-beta` from parent `api/many` setup.
- `UpdateSkill.Name` is the CLI alias shown on status lines.

## Steps

1. Leaves choose which skills to pre-install and whether to mutate disk before update.

## Context

- Shared flag args apply to every skill in the batch.
- Exit remains success when some skills are not installed.
