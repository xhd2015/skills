# Scenario

**Feature**: batch update skips skills that were never installed

```
pre-install skill-alpha only -> HandleUpdateMany
  -> skill-alpha  up to date
  -> skill-beta  not installed
  -> summary; no beta dir created
```

## Preconditions

- `skill-beta` has no `SKILL.md` anywhere.

## Steps

1. Pre-install `skill-alpha` only.
2. Run batch update with default target resolution.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.PreInstalls = []PreInstall{{
		Opts: install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent},
		Args: nil,
	}}
	req.UpdateArgs = nil
	return nil
}
```
