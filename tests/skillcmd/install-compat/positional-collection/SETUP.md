# Scenario

**Feature**: positional `<dir>` with basename `skills` uses the same smart nest as `--dir`

```go
import "github.com/xhd2015/skills/skillcmd"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.InstallOpts = skillcmd.InstallOptions{
		SkillDirName: "demo-skill",
		SkillContent: "# installed via positional\n",
		Usage:        "skill --install",
	}
	req.Args = []string{"vendor/skills"}
	return nil
}
```
