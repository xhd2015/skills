# Scenario

**Feature**: rename nested SKILL.md to TOPIC.md reports delete+create

```
disk: SKILL.md match + a/SKILL.md (same body)
plan: SKILL.md + ExtraFiles a/TOPIC.md (same body)
HandleInstall
  -> create: <abs>/a/TOPIC.md
  -> delete: <abs>/a/SKILL.md
  -> only TOPIC.md remains under a/
```

## Preconditions

- Root matches.
- Disk has nested `a/SKILL.md`; plan wants `a/TOPIC.md` with the same content.

## Steps

1. Pre-seed root + nested SKILL.md.
2. Install with ExtraFiles `a/TOPIC.md`.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	const nestedBody = "# nested topic body\n"
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# test skill\n"},
		{Name: "a/SKILL.md", Content: nestedBody},
	}
	req.SkillContent = "# test skill\n"
	req.ExtraFiles = []install.InstallFile{
		{Path: "a/TOPIC.md", Content: []byte(nestedBody)},
	}
	req.Args = []string{"example-skill"}
	return nil
}
```
