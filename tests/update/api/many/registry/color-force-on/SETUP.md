# Scenario

**Feature**: --color forces ANSI on batch update status tokens even when stdout is a pipe

```
pre-install alpha -> HandleUpdateMany --color
  -> ANSI escapes around status / summary
```

## Steps

1. Pre-install skill-alpha only.
2. Batch update with `--color`.

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
	req.UpdateArgs = []string{"--color"}
	return nil
}
```
