# Scenario

**Feature**: HandleInstall default target is `.agents/skills/<name>`

```
# skillcmd install API (shim target for deprecated install package)
caller -> skillcmd.HandleInstall(opts, nil) -> Installed skill to: .../demo-skill
```

## Steps

1. Configure InstallOpts for demo-skill with simple content.
2. Args empty → default agents path.

```go
import "github.com/xhd2015/skills/skillcmd"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.InstallOpts = skillcmd.InstallOptions{
		SkillDirName: "demo-skill",
		SkillContent: "# installed via skillcmd\n",
	}
	req.Args = nil
	return nil
}
```
