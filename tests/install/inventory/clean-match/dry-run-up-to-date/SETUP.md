# Scenario

**Feature**: dry-run install when disk equals plan prints up to date only

```
example-skill/SKILL.md == plan
HandleInstall(--dry-run example-skill) -> "[dry-run] Skill is up to date: <absDir>"
```

## Preconditions

- Matching SKILL.md only; no orphans.

## Steps

1. Pre-create matching skill dir.
2. Install with `--dry-run`.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# test skill\n"},
	}
	req.SkillContent = "# test skill\n"
	req.Args = []string{"--dry-run", "example-skill"}
	return nil
}
```
