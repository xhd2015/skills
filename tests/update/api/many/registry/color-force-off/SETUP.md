# Scenario

**Feature**: --no-color disables ANSI on batch update even if environment would allow it

```
pre-install alpha -> HandleUpdateMany --no-color -> no ANSI
```

## Steps

1. Pre-install skill-alpha.
2. Batch update with `--no-color`.

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
	}
	req.UpdateArgs = []string{"--no-color"}
	return nil
}
```
