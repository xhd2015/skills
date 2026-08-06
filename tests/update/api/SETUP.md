# Scenario

**Feature**: choose single-skill or registry batch update API

```
# single skill definition
test -> HandleUpdate(opts, args) -> inventory apply on installed targets

# many registry entries share one flag parse
test -> HandleUpdateMany(skills, args) -> per-skill status + trailing summary
```

## Preconditions

- Descendants set `req.UseMany` and populate either `SingleOpts` or `ManySkills`.

## Steps

1. Descendant grouping nodes set the API mode on `Request`.

## Context

- API choice is the highest-impact split in this tree.
