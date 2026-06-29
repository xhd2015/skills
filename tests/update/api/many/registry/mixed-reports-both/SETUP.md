# Scenario

**Feature**: batch update interleaves installed and not-installed stdout lines

```
pre-install skill-alpha only -> HandleUpdateMany -> up-to-date for alpha, not-installed for beta
```

## Preconditions

- Registry order is `skill-alpha` then `skill-beta` by CLI name.

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