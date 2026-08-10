# Scenario

**Feature**: apply install when disk equals plan stays up to date

```
# matching SKILL.md only; no extras, no orphans
example-skill/SKILL.md == plan
HandleInstall(example-skill) -> "Skill is up to date: <absDir>"
```

## Preconditions

- `example-skill/SKILL.md` already matches planned `SkillContent`.
- No extra files on disk.

## Steps

1. Pre-create skill dir with matching SKILL.md.
2. Install with the same content (no ExtraFiles).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# test skill\n"},
	}
	req.SkillContent = "# test skill\n"
	req.Args = []string{"example-skill"}
	return nil
}
```
