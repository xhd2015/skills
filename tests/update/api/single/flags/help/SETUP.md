# Scenario

**Feature**: `--help` on update shows location and dry-run options

```
HandleUpdate --help -> usage including --dry-run and target flags
```

## Preconditions

- None.

## Steps

1. Call update with `["--help"]`.

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
		Usage:        "skills update",
	}
	req.UpdateArgs = []string{"--help"}
	return nil
}
```
