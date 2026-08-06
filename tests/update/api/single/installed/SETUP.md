# Scenario

**Feature**: update when `SKILL.md` already exists

```
HandleInstall -> SKILL.md on disk
HandleUpdate -> inventory compare + polished status
```

## Preconditions

- Leaves pre-install the skill to the default `.agents/skills/<name>` path unless
  noted otherwise.

## Steps

1. Configure `PreInstalls` with matching canonical content before update.

## Context

- Splits on whether on-disk content matches embedded skill content after install.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
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
