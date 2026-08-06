# Scenario

**Feature**: installed skill with matching content stays up to date

```
HandleInstall (canonical) -> HandleUpdate -> "skill-alpha  up to date"
```

## Preconditions

- `skill-alpha` installed to default agents path with canonical content.

## Steps

1. Pre-install via `HandleInstall` with no extra flags.
2. Update with empty args.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.SingleOpts = install.InstallOptions{
		SkillDirName: "skill-alpha",
		SkillContent: skillAlphaContent,
		Usage:        "skill-alpha update",
	}
	req.PreInstalls = []PreInstall{{
		Opts: req.SingleOpts,
		Args: nil,
	}}
	req.UpdateArgs = nil
	return nil
}
```
