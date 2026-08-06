# Scenario

**Feature**: --color and --no-color are mutually exclusive on update

```
HandleUpdate --color --no-color -> error, no skill dirs created
```

## Steps

1. Run single-skill update with both color flags.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.UseMany = false
	req.SingleOpts = install.InstallOptions{
		SkillDirName: "skill-alpha",
		SkillContent: skillAlphaContent,
		Usage:        "skills update",
	}
	req.UpdateArgs = []string{"--color", "--no-color"}
	return nil
}
```
