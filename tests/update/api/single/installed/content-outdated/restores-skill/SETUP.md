# Scenario

**Feature**: update restores canonical `SKILL.md` after drift

```
HandleInstall -> mutate SKILL.md -> HandleUpdate -> Update skill at + canonical bytes
```

## Preconditions

- `skill-alpha` installed under default agents path.

## Steps

1. Pre-install canonical content.
2. Overwrite `SKILL.md` with `# drifted\n`.
3. Update without flags.

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
	req.PostInstallMutate = []FileMutate{{
		RelPath: skillMDPath(skillAgentsDir("skill-alpha")),
		Content: "# drifted\n",
	}}
	req.UpdateArgs = nil
	return nil
}
```