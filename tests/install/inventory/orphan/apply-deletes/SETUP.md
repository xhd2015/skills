# Scenario

**Bug**: apply install deletes unplanned orphan files

```
example-skill/{SKILL.md match, orphan.txt}
HandleInstall(example-skill)
  -> "Update skill at <absDir>"
  -> "  delete: <absDir>/orphan.txt"
  -> orphan.txt removed
```

## Preconditions

- Root SKILL.md matches plan; `orphan.txt` is unplanned.

## Steps

1. Pre-seed matching SKILL.md + orphan.txt.
2. Install with same content, no ExtraFiles.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# test skill\n"},
		{Name: "orphan.txt", Content: "leftover\n"},
	}
	req.SkillContent = "# test skill\n"
	req.Args = []string{"example-skill"}
	return nil
}
```
