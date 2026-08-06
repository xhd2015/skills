# Scenario

**Feature**: batch `--dry-run` reports would-update with planned paths and leaves disk unchanged

```
pre-install + mutate alpha SKILL.md
HandleUpdateMany --dry-run ->
  skill-alpha  would update  (1 update)
    update  <abs>/SKILL.md
  skill-beta  up to date
  summary includes would update and [dry-run]
disk still drifted
```

## Preconditions

- Alpha installed then drifted; beta installed and current.

## Steps

1. Pre-install both skills.
2. Mutate alpha `SKILL.md`.
3. Batch update with `["--dry-run"]`.

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
		{Opts: install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}, Args: nil},
	}
	req.PostInstallMutate = []FileMutate{{
		RelPath: skillMDPath(skillAgentsDir("skill-alpha")),
		Content: "# drifted for dry-run\n",
	}}
	req.UpdateArgs = []string{"--dry-run"}
	return nil
}
```
