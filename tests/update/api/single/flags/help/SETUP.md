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
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	req.SingleOpts = install.InstallOptions{
		SkillDirName: "skill-alpha",
		SkillContent: skillAlphaContent,
		Usage:        "skill-alpha update",
	}
	req.UpdateArgs = []string{"--help"}
	return nil
}
```