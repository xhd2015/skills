# Scenario

**Feature**: skill-dir inventory sync (rsync-like create/update/delete)

```
# plan = {SKILL.md} ∪ ExtraFiles; compare to all regular files under skill dir
HandleInstall -> inventory
  match + no orphans -> "Skill is up to date"
  else -> header + create:/update:/delete: lines (absolute paths)
```

## Preconditions

- Inventory leaves install to an explicit target dir `example-skill` via Args
  (not the default `.agents/skills/...` layout), so abs paths are
  `WorkDir/example-skill/...`.
- Default skill content is `# test skill\n` unless a leaf overrides.
- Nested pre-existing files use slash paths in `PreExistingFile.Name`.

## Steps

1. Set default SkillDirName / SkillContent for inventory leaves.
2. Leaves set PreExisting*, ExtraFiles, Args (including optional `--dry-run`).
3. Assert inspects stdout inventory lines and on-disk files under WorkDir.

## Context

- Stdout detail lines use absolute paths: `create: `, `update: `, `delete: `
  (one space after colon).
- Action order: all creates (sorted by relative path), then updates, then deletes.
- Pure cleanup (planned content matches, only orphans) still uses
  `Update skill at` + only `delete:` lines.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.SkillDirName == "" {
		req.SkillDirName = "example-skill"
	}
	if req.SkillContent == "" {
		req.SkillContent = "# test skill\n"
	}
	return nil
}

// absSkillDir returns the same absolute path install.HandleInstall prints for
// skill dir `name` under workDir (macOS resolves /var → /private/var).
func absSkillDir(t *testing.T, workDir, name string) string {
	t.Helper()
	// Prefer EvalSymlinks on workDir so /var/folders → /private/var/folders.
	base := workDir
	if resolved, err := filepath.EvalSymlinks(workDir); err == nil {
		base = resolved
	}
	skillDir := filepath.Join(base, name)
	if resolved, err := filepath.EvalSymlinks(skillDir); err == nil {
		return resolved
	}
	// Skill dir may be created during Run; Abs relative to workDir matches install.
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(base); err != nil {
		t.Fatalf("chdir workdir: %v", err)
	}
	abs, absErr := filepath.Abs(name)
	_ = os.Chdir(prev)
	if absErr != nil {
		t.Fatalf("abs skill dir: %v", absErr)
	}
	return abs
}
```
