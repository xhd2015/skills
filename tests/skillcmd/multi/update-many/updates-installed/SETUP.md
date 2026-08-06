# Scenario

**Feature**: update refreshes an already-installed skill when content drifted

```
# pre-install foo, mutate SKILL.md, then update
caller -> HandleSkills(update) -> foo  updated (+ file lines) / restored content
```

## Steps

1. Pre-install foo with canonical content.
2. Mutate on-disk SKILL.md to stale content.
3. Run `skills update`.

```go
import (
	"github.com/xhd2015/doctest/session"
	"path/filepath"

	"github.com/xhd2015/skills/skillcmd"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreInstall = true
	req.PreInstallOpts = skillcmd.InstallOptions{
		SkillDirName: "foo",
		SkillContent: fooSkillContent,
	}
	req.PreInstallArgs = nil
	req.PostInstallMutate = map[string]string{
		filepath.ToSlash(filepath.Join(".agents", "skills", "foo", "SKILL.md")): "# stale foo\n",
	}
	return nil
}
```
