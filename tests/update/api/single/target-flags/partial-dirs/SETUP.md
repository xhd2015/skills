# Scenario

**Feature**: partial install across multiple flag targets

```
install --codex only -> update --codex --opencode -> only codex processed
```

## Preconditions

- `skill-alpha` installed with `--codex` only.

## Steps

1. Pre-install with `["--codex"]`.
2. Drift codex `SKILL.md`.
3. Update with `["--codex", "--opencode"]`.

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