# Scenario

**Feature**: `--global` update affects home dir only

```
HandleInstall --global -> ~/.agents/skills/<name>/SKILL.md
HandleUpdate --global -> refresh global copy; local project dir untouched
```

## Preconditions

- `req.UseGlobalHome` is true so `HOME` points at a temp directory.

## Steps

1. Pre-install with `["--global"]`.
2. Drift global `SKILL.md`.
3. Update with `["--global"]`.

```go
import (
	"path/filepath"

	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, req *Request) error {
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