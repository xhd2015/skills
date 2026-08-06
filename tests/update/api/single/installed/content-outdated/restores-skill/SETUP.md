# Scenario

**Feature**: update restores canonical `SKILL.md` after drift

```
HandleInstall -> mutate SKILL.md -> HandleUpdate
  -> skill-alpha  updated  (1 update)
  ->   update  <abs>/SKILL.md
  -> canonical bytes on disk
```

## Preconditions

- `skill-alpha` installed under default agents path.

## Steps

1. Pre-install canonical content.
2. Overwrite `SKILL.md` with `# drifted\n`.
3. Update without flags.

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
	req.PostInstallMutate = []FileMutate{{
		RelPath: skillMDPath(skillAgentsDir("skill-alpha")),
		Content: "# drifted\n",
	}}
	req.UpdateArgs = nil
	return nil
}
```
