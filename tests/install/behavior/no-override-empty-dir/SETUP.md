## Preconditions
- An empty directory "example-skill" exists (no files inside, no SKILL.md).
- The `--no-override` flag is passed.

## Steps
1. Create an empty pre-existing directory (no PreExistingFiles).
2. Call `HandleInstall` with `--no-override` and fresh skill content `"# new skill\n"`.

## Context
- `--no-override` should only abort when the target directory is **non-empty**.
- An empty directory should be treated as a fresh install.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingDir = "example-skill"
	// No PreExistingFiles — directory is empty
	req.Args = []string{"--no-override", "example-skill"}
	req.SkillContent = "# new skill\n"
	return nil
}
```
