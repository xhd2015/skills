# Scenario

**Feature**: `--global` update affects home dir only

```
HandleInstall --global -> ~/.agents/skills/<name>/SKILL.md
HandleUpdate --global
  -> skill-alpha  updated  (1 update)
  ->   update  <abs-home>/.agents/skills/skill-alpha/SKILL.md
  -> local project dir untouched
```

## Preconditions

- `req.UseGlobalHome` is true so `HOME` points at a temp directory.

## Steps

1. Pre-install with `["--global"]`.
2. Drift global `SKILL.md` via `$HOME/` mutate path.
3. Update with `["--global"]`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.UseGlobalHome = true
	req.SingleOpts = install.InstallOptions{
		SkillDirName: "skill-alpha",
		SkillContent: skillAlphaContent,
		Usage:        "skill-alpha update",
	}
	req.PreInstalls = []PreInstall{{
		Opts: req.SingleOpts,
		Args: []string{"--global"},
	}}
	globalSkill := "$HOME/" + filepath.ToSlash(filepath.Join(skillAgentsDir("skill-alpha"), "SKILL.md"))
	req.PostInstallMutate = []FileMutate{{
		RelPath: globalSkill,
		Content: "# drifted global\n",
	}}
	req.UpdateArgs = []string{"--global"}
	return nil
}
```
