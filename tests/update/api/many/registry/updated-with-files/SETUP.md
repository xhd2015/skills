# Scenario

**Feature**: batch update prints indented absolute file paths when content changes

```
pre-install alpha (SKILL.md only) + beta
mutate alpha SKILL.md; registry alpha plans extra.md
HandleUpdateMany ->
  skill-alpha  updated  (1 create, 1 update)
    create  <abs>/extra.md
    update  <abs>/SKILL.md
  skill-beta  up to date
  summary 1 updated · 1 up to date · 0 not installed
```

## Preconditions

- Alpha and beta installed with SKILL.md only.
- Alpha on-disk content drifted; update plan also adds `extra.md` via ExtraFiles.

## Steps

1. Pre-install both skills without extras.
2. Mutate alpha `SKILL.md`.
3. Configure `ManySkills` so alpha includes ExtraFiles `extra.md`.
4. Batch update with no flags.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	alphaInstall := install.InstallOptions{SkillDirName: "skill-alpha", SkillContent: skillAlphaContent}
	betaInstall := install.InstallOptions{SkillDirName: "skill-beta", SkillContent: skillBetaContent}
	req.PreInstalls = []PreInstall{
		{Opts: alphaInstall, Args: nil},
		{Opts: betaInstall, Args: nil},
	}
	req.PostInstallMutate = []FileMutate{{
		RelPath: skillMDPath(skillAgentsDir("skill-alpha")),
		Content: "# drifted alpha\n",
	}}
	req.ManySkills = []install.UpdateSkill{
		{
			Name: "skill-alpha",
			InstallOptions: install.InstallOptions{
				SkillDirName: "skill-alpha",
				SkillContent: skillAlphaContent,
				ExtraFiles: []install.InstallFile{
					{Path: "extra.md", Content: []byte(extraFileContent)},
				},
			},
		},
		{
			Name:           "skill-beta",
			InstallOptions: betaInstall,
		},
	}
	req.UpdateArgs = nil
	return nil
}
```
