# Scenario

**Feature**: `HandleUpdateMany` walks a skill registry

```
test harness -> HandleUpdateMany([]UpdateSkill, args) -> each skill's installed targets
```

## Preconditions

- `req.UseMany` is true for all leaves under this branch.

## Steps

1. Set `req.UseMany = true`.
2. Leaves build `req.ManySkills` and `req.UpdateArgs`.

## Context

- Registry leaves use `skill-alpha` and `skill-beta` with distinct canonical content.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	req.UseMany = true
	if len(req.ManySkills) == 0 {
		req.ManySkills = []install.UpdateSkill{
			{InstallOptions: install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}},
			{InstallOptions: install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}},
		}
	}
	return nil
}
```