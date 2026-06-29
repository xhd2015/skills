# Scenario

**Feature**: installed skill with matching content stays up to date

```
HandleInstall (canonical) -> HandleUpdate -> "Skill is up to date"
```

## Preconditions

- `skill-alpha` installed to default agents path with canonical content.

## Steps

1. Pre-install via `HandleInstall` with no extra flags.
2. Update with empty args.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
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