# Scenario

**Feature**: fresh install when skill directory is missing

```
# no example-skill dir on disk
HandleInstall -> "Installed skill to: <absDir>" + create: lines for plan files
```

## Preconditions

- Target skill directory does not exist before install.

## Steps

1. Leaves set plan (SkillContent + ExtraFiles) without PreExistingDir.
2. Install to explicit `example-skill` path.

## Context

- Fresh installs only emit `create:` detail lines (sorted by relative path).
- Header is always `Installed skill to:` when the dir was missing.

```go
func Setup(t *testing.T, req *Request) error {
	// Fresh leaves must not pre-create the skill dir.
	req.PreExistingDir = ""
	req.PreExistingFiles = nil
	req.Args = []string{"example-skill"}
	return nil
}
```
