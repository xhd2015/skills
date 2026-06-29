# Scenario

**Feature**: batch update reports each installed skill as up to date

```
pre-install alpha + beta -> HandleUpdateMany -> two up-to-date lines
```

## Preconditions

- Both skills installed to default agents paths with canonical content.

## Steps

1. Pre-install both skills.
2. Batch update with no flags.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	req.PreInstalls = []PreInstall{
		{Opts: install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}, Args: nil},
		{Opts: install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}, Args: nil},
	}
	req.UpdateArgs = nil
	return nil
}
```