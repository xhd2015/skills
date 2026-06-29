# Scenario

**Feature**: update handlers refresh installed skills only

```
# optional pre-install seeds SKILL.md via HandleInstall
test harness -> HandleInstall -> target dirs with SKILL.md

# update resolves targets and skips dirs without SKILL.md
test harness -> HandleUpdate / HandleUpdateMany -> InstallTo per installed dir
test harness <- stdout (up to date / update / dry-run messages)
```

## Preconditions

- `HandleUpdate`, `HandleUpdateMany`, and `HandleInstall` are available from
  `github.com/xhd2015/skills/install`.
- Each test runs in an isolated temporary working directory.

## Steps

1. Create a temp directory and `chdir` into it.
2. Optionally set `HOME` when `req.UseGlobalHome` is true.
3. Run `req.PreInstalls` via `HandleInstall` (stdout not captured).
4. Apply `req.PostInstallMutate` writes to simulate drifted on-disk content.
5. Capture stdout while calling `HandleUpdate` or `HandleUpdateMany`.
6. Return stdout, error string, and `WorkDir` in the response.

## Context

- Update flag surface: `--global`, `--cursor`, `--codex`, `--opencode`,
  `--general-agents`, `--dry-run`, `-h` / `--help`.
- Install-only flags (`--force`, `--no-override`) are not accepted on update.
- Canonical test skill names: `skill-alpha`, `skill-beta` (used by registry leaves).

```go
import (
	"os"
	"path/filepath"
	"testing"
)

const (
	skillAlphaContent = "# skill alpha canonical\n"
	skillBetaContent  = "# skill beta canonical\n"
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

func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```