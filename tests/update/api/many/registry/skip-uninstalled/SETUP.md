# Scenario

**Feature**: batch update ignores skills that were never installed

```
pre-install skill-alpha only -> HandleUpdateMany -> stdout only for alpha
```

## Preconditions

- `skill-beta` has no `SKILL.md` anywhere.

## Steps

1. Pre-install `skill-alpha` only.
2. Run batch update with default target resolution.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	alpha := install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}
	req.PreInstalls = []PreInstall{{
		Opts: alpha,
		Args: nil,
	}}
	req.UpdateArgs = nil
	return nil
}
```