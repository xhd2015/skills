# Scenario

**Feature**: update handlers refresh installed skills only

```
# parent suite stays Parallel-safe; child owns cwd / HOME / product stdout
test suite Run -> cmd/runupdate (child process, WorkDir + optional HOME)
child -> HandleInstall (stdout discarded) -> mutate disk
child -> HandleUpdate / HandleUpdateMany -> polished stdout
test suite <- captured stdout + WorkDir (+ HomeDir when global)
```

## Preconditions

- `HandleUpdate`, `HandleUpdateMany`, and `HandleInstall` are available from
  `github.com/xhd2015/skills/install`.
- Each leaf gets an isolated temporary `WorkDir` (child process chdirs there).
- Suite `Run` must not use `os.Chdir` / `t.Setenv` / `os.Stdout` reassignment.

## Steps

1. Create a temp `WorkDir` (and optional isolated `HomeDir` for `--global`).
2. Serialize the request for `cmd/runupdate` and `go run` it with `GOWORK=off`.
3. Child: chdir, optional HOME, pre-install (stdout discarded), mutate, update.
4. Return captured product stdout, error string, `WorkDir`, and `HomeDir`.

## Context

- Update flag surface: `--global`, `--cursor`, `--codex`, `--opencode`,
  `--general-agents`, `--dry-run`, `-h` / `--help`.
- Install-only flags (`--force`, `--no-override`) are not accepted on update.
- Canonical test skill names: `skill-alpha`, `skill-beta` (used by registry leaves).
- Status line spacing: `{name}  {status}` (two spaces). File lines:
  `  {op}  {absPath}` (two-space indent, two spaces after op, no colon).
- Batch always ends with a blank line + summary + trailing newline.

```go
import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

const (
	skillAlphaContent = "# skill alpha canonical\n"
	skillBetaContent  = "# skill beta canonical\n"
	extraFileContent  = "# extra file\n"
)

func skillAgentsDir(skillName string) string {
	return filepath.Join(".agents", "skills", skillName)
}

func skillCodexDir(skillName string) string {
	return filepath.Join(".codex", "skills", skillName)
}

func skillOpencodeDir(skillName string) string {
	return filepath.Join(".opencode", "skills", skillName)
}

func skillMDPath(dir string) string {
	return filepath.Join(dir, "SKILL.md")
}

func absUnder(workDir string, rel string) string {
	return filepath.Join(workDir, rel)
}

func rePath(p string) string {
	return regexp.QuoteMeta(p)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func pathExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	return err == nil
}

// assertBatchStdoutPolished checks leading/trailing newline polish.
func assertBatchStdoutPolished(t *testing.T, stdout string) {
	t.Helper()
	if strings.HasPrefix(stdout, "\n") {
		t.Fatalf("stdout must not start with a blank line: %q", stdout)
	}
	if stdout != "" && !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout must end with newline: %q", stdout)
	}
}

// assertNoLegacyUpdateStdout rejects pre-polish InstallTo / batch strings.
func assertNoLegacyUpdateStdout(t *testing.T, stdout string) {
	t.Helper()
	legacy := []string{
		"Skill is up to date:",
		"skill not installed:",
		"Update skill at ",
		"[dry-run] Skill is up to date:",
		"[dry-run] Update skill at",
	}
	for _, frag := range legacy {
		if strings.Contains(stdout, frag) {
			t.Fatalf("legacy stdout fragment %q still present:\n%s", frag, stdout)
		}
	}
	// Bare trailing skill-name-only lines after inventory were old batch markers.
	// New format never prints a line that is only the skill dir name.
}
```
