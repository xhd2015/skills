# Scenario

**Feature**: ExtraFiles path validation rejects illegal paths

```
# invalid ExtraFiles path short-circuits before inventory writes
HandleInstall(ExtraFiles with invalid Path) -> error; no skill dir written
```

## Preconditions
- We are testing extra file path validation — invalid paths should produce errors before any directory modification.
- These tests verify that `resolveInstallFiles` correctly rejects invalid paths.

## Context
- Extra file paths are validated early in `installTo` → `resolveInstallFiles`.
- Invalid paths include: ".", "..", absolute paths, paths starting with ".."+os.PathSeparator, and "SKILL.md".
- All invalid paths produce an error; no directory or file is created.
- The default skill dir name is "example-skill" and default content is "# test skill\n".

```go
func Setup(t *testing.T, req *Request) error {
	if req.SkillDirName == "" {
		req.SkillDirName = "example-skill"
	}
	if req.SkillContent == "" {
		req.SkillContent = "# test skill\n"
	}
	return nil
}
```
