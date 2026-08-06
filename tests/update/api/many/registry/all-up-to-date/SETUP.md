# Scenario

**Feature**: batch update reports each installed skill as up to date

```
pre-install alpha + beta -> HandleUpdateMany
  -> skill-alpha  up to date
  -> skill-beta  up to date
  -> summary 0 updated · 2 up to date · 0 not installed
```

## Preconditions

- Both skills installed to default agents paths with canonical content.

## Steps

1. Pre-install both skills.
2. Batch update with no flags.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.PreInstalls = []PreInstall{
		{Opts: install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}, Args: nil},
		{Opts: install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}, Args: nil},
	}
	req.UpdateArgs = nil
	return nil
}
```
