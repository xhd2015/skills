# Scenario

**Feature**: fresh install creates root and nested TOPIC.md extras

```
no example-skill/
plan: SKILL.md + nested/TOPIC.md
HandleInstall
  -> Installed skill to: <absDir>
  -> create: <absDir>/SKILL.md
  -> create: <absDir>/nested/TOPIC.md
```

## Preconditions

- Target dir missing.
- Plan includes one nested ExtraFile `nested/TOPIC.md`.

## Steps

1. Set ExtraFiles and matching SkillContent.
2. Install to `example-skill`.

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SkillContent = "# test skill\n"
	req.ExtraFiles = []install.InstallFile{
		{Path: "nested/TOPIC.md", Content: []byte("# nested topic\n")},
	}
	req.Args = []string{"example-skill"}
	return nil
}
```
