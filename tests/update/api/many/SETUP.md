# Scenario

**Feature**: `HandleUpdateMany` walks a skill registry

```
test harness -> HandleUpdateMany([]UpdateSkill, args)
  -> per skill status line (+ indented file ops when changed)
  -> trailing summary counts
```

## Preconditions

- `req.UseMany` is true for all leaves under this branch.

## Steps

1. Set `req.UseMany = true`.
2. Leaves build `req.ManySkills` and `req.UpdateArgs`.

## Context

- Registry leaves use `skill-alpha` and `skill-beta` with distinct canonical content.
- Batch stdout always includes a trailing summary line.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.UseMany = true
	if len(req.ManySkills) == 0 {
		req.ManySkills = []install.UpdateSkill{
			{InstallOptions: install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}, Name: "skill-alpha"},
			{InstallOptions: install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}, Name: "skill-beta"},
		}
	}
	return nil
}
```
