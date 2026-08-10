# Scenario

**Feature**: `--no-override` on an empty pre-existing dir does not abort

```
# empty example-skill/ exists; plan needs create for SKILL.md
HandleInstall(--no-override example-skill)
  -> no confirmation / no Aborted
  -> Update skill at (dir already existed)
```

## Preconditions
- An empty directory "example-skill" exists (no files inside, no SKILL.md).
- The `--no-override` flag is passed.

## Steps
1. Create an empty pre-existing directory (no PreExistingFiles).
2. Call `HandleInstall` with `--no-override` and fresh skill content `"# new skill\n"`.

## Context
- `--no-override` should only require confirmation when the target directory is
  **non-empty** and the plan needs create/update/delete.
- An empty directory needs no confirmation; header is still `Update skill at`
  because the directory already existed.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreExistingDir = "example-skill"
	// No PreExistingFiles — directory is empty
	req.Args = []string{"--no-override", "example-skill"}
	req.SkillContent = "# new skill\n"
	return nil
}
```
