# Scenario

**Feature**: differing planned content updates the existing file

```
example-skill/SKILL.md = "# old skill\n"
plan SKILL.md = "# new skill\n"
HandleInstall -> update: <absDir>/SKILL.md; content becomes new
```

## Preconditions

- Skill dir exists with old SKILL.md content; plan has new content.
- No ExtraFiles; no orphans.

## Steps

1. Pre-seed old SKILL.md.
2. Install with new SkillContent.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# old skill\n"},
	}
	req.SkillContent = "# new skill\n"
	req.Args = []string{"example-skill"}
	return nil
}
```
