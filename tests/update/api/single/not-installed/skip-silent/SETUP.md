# Scenario

**Feature**: update with no prior install produces no output

```
# no SKILL.md at resolved targets
HandleUpdate -> (skip each dir silently) -> no stdout
```

## Preconditions

- No `PreInstalls`; default target `.agents/skills/skill-alpha` does not exist.

## Steps

1. Call update with default args (no location flags).

## Context

- Confirms update never creates a fresh install.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	req.SingleOpts = install.InstallOptions{
		SkillDirName: "skill-alpha",
		SkillContent: skillAlphaContent,
		Usage:        "skill-alpha update",
	}
	req.UpdateArgs = nil
	return nil
}
```