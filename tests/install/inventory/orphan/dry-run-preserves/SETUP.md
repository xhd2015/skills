# Scenario

**Bug**: dry-run reports orphan delete without removing the file

```
example-skill/{SKILL.md match, orphan.txt}
HandleInstall(--dry-run example-skill)
  -> "[dry-run] Update skill at <absDir>"
  -> "[dry-run]   delete: <absDir>/orphan.txt"
  -> orphan.txt still on disk
```

## Preconditions

- Matching SKILL.md + unplanned orphan.txt.

## Steps

1. Pre-seed files.
2. Install with `--dry-run`.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# test skill\n"},
		{Name: "orphan.txt", Content: "leftover\n"},
	}
	req.SkillContent = "# test skill\n"
	req.Args = []string{"--dry-run", "example-skill"}
	return nil
}
```
