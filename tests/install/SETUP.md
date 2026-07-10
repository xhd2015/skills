# Scenario

**Feature**: HandleInstall installs and inventory-syncs skill directories

```
# caller installs skill content + extras into a target skill dir
caller -> install.HandleInstall(opts, args) -> skill dir under WorkDir

# inventory compares plan {SKILL.md}∪ExtraFiles to all regular files under skill dir
HandleInstall -> inventory sync
  -> create: / update: / delete:  (or Skill is up to date)

# dry-run prints the same plan with [dry-run] prefix and does not write
caller --dry-run -> HandleInstall -> stdout only
```

## Preconditions
- The install package's `HandleInstall` function is available at the module root
  (shim over `skillcmd.HandleInstall`).
- Tests run in an isolated temporary directory.
- Working directory is changed to a fresh temp dir before each test.
- Pre-existing nested files under `PreExistingDir` may use slash paths; Run
  creates parent directories via `MkdirAll`.

## Steps
1. Create a temporary directory and change into it.
2. If `req.UseGlobalHome` is true, set `HOME` to a separate temporary directory.
3. If `req.PreExistingDir` is non-empty, create the directory with the specified
   pre-existing files (nested `Name` values get parent dirs).
4. If `req.NonInteractive` is true, replace stdin with an empty file so
   `confirmOverwrite` returns false without waiting for user input.
5. Capture stdout via an os.Pipe.
6. Call `HandleInstall(InstallOptions{SkillDirName, SkillContent, ExtraFiles, Usage: ""}, req.Args)`.
7. Collect captured stdout and any error returned by HandleInstall into the Response.

## Context
- `HandleInstall` is the entry point that parses flags from `args` and calls
  inventory sync for each target directory.
- Supported flags: `--cursor`, `--codex`, `--opencode`, `--general-agents`,
  `--global`, `--no-override`, `--force`, `--dry-run`.
- Output messages (inventory era):
  - Up-to-date (plan match + no orphans): `"Skill is up to date: <absDir>\n"`
  - Fresh install (skill dir missing): `"Installed skill to: <absDir>\n"` then
    sorted `create: <absPath>` lines
  - Update (dir existed): `"Update skill at <absDir>\n"` then sorted
    `create:` / `update:` / `delete:` lines (absolute paths)
  - Dry-run variants: same lines prefixed with `"[dry-run] "`
  - Aborted (user declined confirmation): `"Aborted."`
- Detail path convention for inventory leaves: **absolute** paths under the skill dir.
- Extra file paths are validated in `resolveInstallFiles`: ".", "..", absolute
  paths, and "SKILL.md" are rejected with `"invalid install file path"` or
  `"extra install file cannot replace SKILL.md"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Normalize slices so leaf Setups can append without nil checks.
	if req.Args == nil {
		req.Args = []string{}
	}
	if req.ExtraFiles == nil {
		req.ExtraFiles = nil // keep nil semantics for "unset" vs empty slice
	}
	if req.PreExistingFiles == nil {
		req.PreExistingFiles = []PreExistingFile{}
	}
	return nil
}
```
