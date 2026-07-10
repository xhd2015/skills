# Scenario

**Feature**: existing skill dir with different content prints Update header

```
example-skill/{SKILL.md old, stale.txt}
HandleInstall(new content) -> "Update skill at <absDir>"
```

## Preconditions
- A directory "example-skill" exists in the working directory.
- It contains a SKILL.md with content `"# old skill\n"` (different from the new content `"# new skill\n"`).
- It also contains a stale file "stale.txt" (orphan under inventory sync).

## Steps
1. Create the pre-existing directory with old skill content and a stale file.
2. Call `HandleInstall` with new skill content different from the old.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# old skill\n"},
		{Name: "stale.txt", Content: "stale"},
	}
	req.Args = []string{"example-skill"}
	req.SkillContent = "# new skill\n"
	return nil
}
```
