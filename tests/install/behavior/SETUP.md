# Scenario

**Feature**: HandleInstall behavior flags and basic install headers

```
# fresh vs existing target dirs and flag interactions
HandleInstall(args) -> Installed skill to: | Update skill at | Aborted.
```

## Preconditions
- We are testing normal install behavior — stdout output messages and filesystem side effects.
- The focus is on what `HandleInstall` prints and what files it creates, not on error paths.

## Context
- The default skill directory name is "example-skill".
- The default skill content is "# test skill\n".
- These defaults can be overridden by leaf Setup functions.

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
