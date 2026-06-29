# Scenario

**Feature**: `HandleUpdate` for one embedded skill

```
test harness -> HandleUpdate(InstallOptions, updateArgs) -> InstallTo on installed targets only
```

## Preconditions

- `req.UseMany` is false for all leaves under this branch.

## Steps

1. Set `req.UseMany = false`.
2. Leaves configure `req.SingleOpts` and `req.UpdateArgs`.

## Context

- Default target when no location flag is passed matches install:
  `.agents/skills/<SkillDirName>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.UseMany = false
	return nil
}
```