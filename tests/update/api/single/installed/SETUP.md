# Scenario

**Feature**: update when `SKILL.md` already exists

```
HandleInstall -> SKILL.md on disk
HandleUpdate -> InstallTo compares content
```

## Preconditions

- Leaves pre-install the skill to the default `.agents/skills/<name>` path unless
  noted otherwise.

## Steps

1. Configure `PreInstalls` with matching canonical content before update.

## Context

- Splits on whether on-disk content matches embedded skill content after install.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	if req.SingleOpts.SkillDirName == "" {
		req.SingleOpts = install.InstallOptions{
			SkillDirName: "skill-alpha",
			SkillContent: skillAlphaContent,
			Usage:        "skill-alpha update",
		}
	}
	return nil
}
```