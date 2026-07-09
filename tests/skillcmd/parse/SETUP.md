# Scenario

**Feature**: ParseSkillArgs classifies skill action flags and rest args

```
# scan argv for --show / --install / --list / --header
caller -> skillcmd.ParseSkillArgs(args) -> Action + Header + Rest | error
```

## Preconditions

- Mode is parse; no filesystem side effects.

## Steps

1. Set `req.Mode = ModeParse`.
2. Leaves set `req.Args` for each parse case.

## Context

- Install option flags such as `--global` remain in Rest.
- Exactly one of show / install / list is required.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = ModeParse
	return nil
}
```
