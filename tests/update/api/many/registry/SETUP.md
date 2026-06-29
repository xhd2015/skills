# Scenario

**Feature**: registry batch reports not-installed skills and updates installed ones

```
HandleUpdateMany -> per skill: InstallTo if any SKILL.md else skill not installed line
```

## Preconditions

- `ManySkills` defaults to `skill-alpha` and `skill-beta` from parent `api/many` setup.
- Implementer sets `UpdateSkill.Name` to CLI alias; here `SkillDirName` matches alias.

## Steps

1. Leaves choose which skills to pre-install before batch update.

## Context

- Shared flag args apply to every skill in the batch.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	if len(req.ManySkills) == 0 {
		req.ManySkills = []install.UpdateSkill{
			{InstallOptions: install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}},
			{InstallOptions: install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}},
		}
	}
	return nil
}
```