## Preconditions
- The install package's `HandleInstall` function is available at the module root.
- Tests run in an isolated temporary directory.
- Working directory is changed to a fresh temp dir before each test.

## Steps
1. Create a temporary directory and change into it.
2. If `req.UseGlobalHome` is true, set `HOME` to a separate temporary directory.
3. If `req.PreExistingDir` is non-empty, create the directory with the specified pre-existing files.
4. If `req.NonInteractive` is true, replace stdin with an empty file so `confirmOverwrite` returns false without waiting for user input.
5. Capture stdout via an os.Pipe.
6. Call `HandleInstall(InstallOptions{SkillDirName, SkillContent, ExtraFiles, Usage: ""}, req.Args)`.
7. Collect captured stdout and any error returned by HandleInstall into the Response.

## Context
- `HandleInstall` is the entry point that parses flags from `args` and calls `installTo` for each target directory.
- Supported flags: `--cursor`, `--codex`, `--opencode`, `--general-agents`, `--global`, `--no-override`, `--force`, `--dry-run`.
- Output messages:
  - Fresh install (directory does not exist): `"Installed skill to: <absDir>"`
  - Overwrite (directory exists and is replaced): `"Update skill at <absDir>"`
  - Up-to-date (SKILL.md unchanged): `"Skill is up to date: <absDir>"`
  - Dry-run variants: prefixed with `"[dry-run] "`
  - Aborted (user declined confirmation): `"Aborted."`
- Extra file paths are validated in `resolveInstallFiles`: ".", "..", absolute paths, and "SKILL.md" are rejected with `"invalid install file path"` or `"extra install file cannot replace SKILL.md"`.

