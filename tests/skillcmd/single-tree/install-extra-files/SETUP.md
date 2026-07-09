# Scenario

**Feature**: install writes nested skill-cli/SKILL.md extras (not topics/*.md)

```
caller -> SingleSkill.Handle(--install) -> .agents/.../skill-cli/SKILL.md
```

## Preconditions

- Tree SingleSkill configured by parent.

## Steps

1. Configure Args (and ExtraFiles when installing).

```go
import "github.com/xhd2015/skills/skillcmd"

func Setup(t *testing.T, req *Request) error {
	req.UseWorkDir = true
	req.Args = []string{"--install"}
	req.ExtraFiles = []skillcmd.InstallFile{
		{Path: "skill-cli/SKILL.md", Content: []byte(nestedSkillCLIContent)},
	}
	return nil
}
```
