# Scenario

**Feature**: nested file removed from plan is deleted on install

```
disk: SKILL.md match + a/TOPIC.md
plan: SKILL.md only (no ExtraFiles)
HandleInstall -> delete: <absDir>/a/TOPIC.md; nested file gone
```

## Preconditions

- Root matches plan; nested `a/TOPIC.md` exists on disk but is not planned.

## Steps

1. Pre-seed root + nested topic file.
2. Install with matching root content and empty ExtraFiles.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# test skill\n"},
		{Name: "a/TOPIC.md", Content: "# nested topic\n"},
	}
	req.SkillContent = "# test skill\n"
	// ExtraFiles intentionally empty / nil — nested path is no longer planned.
	req.Args = []string{"example-skill"}
	return nil
}
```
