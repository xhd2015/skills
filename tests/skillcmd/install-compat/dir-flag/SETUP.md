# Scenario

**Feature**: `--dir` selects an explicit install destination with smart layout

```
HandleInstall(--dir …) -> ResolveExplicitSkillDir -> skill root
```

## Preconditions

- Parent install-compat setup (ModeInstallCompat, UseWorkDir).
- Leaves set `InstallOpts` and `--dir` args.

```go
import "github.com/xhd2015/skills/skillcmd"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.InstallOpts = skillcmd.InstallOptions{
		SkillDirName: "demo-skill",
		SkillContent: "# installed via --dir\n",
		Usage:        "skill --install",
	}
	return nil
}
```
