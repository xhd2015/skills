# Scenario

**Feature**: `--dry-run` previews update without writing drifted file back

```
mutate SKILL.md -> HandleUpdate --dry-run
  -> skill-alpha  would update  (1 update)
  ->   update  <abs>/SKILL.md
  -> drift remains
```

## Preconditions

- Same pre-install and drift as `restores-skill`.

## Steps

1. Pre-install and mutate `SKILL.md`.
2. Pass `--dry-run` to update.

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
	drifted := "# drifted for dry-run\n"
	req.PostInstallMutate = []FileMutate{{
		RelPath: skillMDPath(skillAgentsDir("skill-alpha")),
		Content: drifted,
	}}
	req.UpdateArgs = []string{"--dry-run"}
	return nil
}
```
