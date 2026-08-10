# Scenario

**Feature**: `--force` overrides `--no-override` without confirmation

```
example-skill/ non-empty
HandleInstall(--force --no-override) -> Update skill at; no Aborted
```

## Preconditions
- A non-empty directory "example-skill" exists with a SKILL.md containing `"# old skill\n"`.
- Both `--force` and `--no-override` flags are passed.

## Steps
1. Create the pre-existing directory with old skill content.
2. Call `HandleInstall` with `--force --no-override` and new skill content `"# new skill\n"`.

## Context
- `--force` takes precedence over `--no-override` (the code sets `noOverride = false` when `force` is true).
- The directory should be updated without triggering the confirmation prompt.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreExistingDir = "example-skill"
	req.PreExistingFiles = []PreExistingFile{
		{Name: "SKILL.md", Content: "# old skill\n"},
	}
	req.Args = []string{"--force", "--no-override", "example-skill"}
	req.SkillContent = "# new skill\n"
	return nil
}
```
