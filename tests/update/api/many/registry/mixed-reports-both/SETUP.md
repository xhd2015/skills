# Scenario

**Feature**: batch update interleaves installed and not-installed status lines

```
pre-install skill-alpha only -> HandleUpdateMany
  -> registry order: alpha up to date, then beta not installed
  -> summary counts both
```

## Preconditions

- Registry order is `skill-alpha` then `skill-beta` by CLI name.

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
