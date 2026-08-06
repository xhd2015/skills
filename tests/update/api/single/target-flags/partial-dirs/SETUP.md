# Scenario

**Feature**: partial install across multiple flag targets

```
install --codex only -> update --codex --opencode
  -> only codex processed (updated status + file line)
  -> opencode still absent
```

## Preconditions

- `skill-alpha` installed with `--codex` only.

## Steps

1. Pre-install with `["--codex"]`.
2. Drift codex `SKILL.md`.
3. Update with `["--codex", "--opencode"]`.

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
		Args: []string{"--codex"},
	}}
	codexSkill := skillMDPath(skillCodexDir("skill-alpha"))
	req.PostInstallMutate = []FileMutate{{
		RelPath: codexSkill,
		Content: "# drifted codex\n",
	}}
	req.UpdateArgs = []string{"--codex", "--opencode"}
	return nil
}
```
