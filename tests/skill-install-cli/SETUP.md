# Scenario

**Feature**: repo skill CLIs support `skill --install` with dry-run and global targets

```
# user invokes skill --install on a CLI binary
user -> skill CLI (skill --install) -> install/skillcmd HandleInstall -> stdout / filesystem

# global scope resolves under HOME
user -> skill CLI (skill --install --global) -> ~/.agents/skills/<name>

# go-best-practice ships nested cli/skill-cli/TOPIC.md extras (Shape 3)
HandleInstall -> SKILL.md + cli/skill-cli/TOPIC.md (+ other nested paths)
```

## Preconditions

- Module root contains `go.mod` for `github.com/xhd2015/skills`.
- All three CLI binaries are built once per `doctest test` session into a temp
  cache keyed by `DOCTEST_SESSION_ID`.
- Leaves set `req.Binary` via CLI grouping setup or inherit from ancestor setup.

## Steps

1. Build `go-best-practice`, `playwright-debug`, and `github-fetch` via session cache.
2. Each leaf sets `req.Args` and scope flags (`UseWorkDir`, `UseGlobalHome`).
3. `Run` executes the binary and captures stdout, stderr, and exit code.

## Context

- Tests shell out to real CLI binaries (not the `install` package directly).
- Local install tests chdir into an isolated temp work directory.
- Global install tests set `HOME` to a separate temp directory.
- After skillcmd migration, action flags replace word subcommands under `skill`.

```go
import (
	"sync"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func moduleRoot(t *testing.T, d *session.Doctest) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	if err != nil {
		t.Fatalf("module root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatalf("module root missing go.mod at %s: %v", root, err)
	}
	return root
}





// Process-local CLI binary cache (one-process; in-memory mutex, not session flock).
var (
	cliBinMu    sync.Mutex
	cliBinPaths = map[string]string{}
	cliBinErrs  = map[string]error{}
)

func buildCLIBinaryOnce(t *testing.T, d *session.Doctest, name, pkg string) (string, error) {
	t.Helper()
	cliBinMu.Lock()
	defer cliBinMu.Unlock()
	if p, ok := cliBinPaths[name]; ok {
		return p, cliBinErrs[name]
	}
	if err, ok := cliBinErrs[name]; ok && err != nil {
		return "", err
	}
	root := moduleRoot(t, d)
	dir, err := os.MkdirTemp("", "skills-cli-"+name+"-")
	if err != nil {
		cliBinErrs[name] = err
		return "", err
	}
	bin := filepath.Join(dir, name)
	cmd := exec.Command("go", "build", "-o", bin, pkg)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		e := fmt.Errorf("build %s: %w\n%s", name, err, out)
		cliBinErrs[name] = e
		return "", e
	}
	if err := os.Chmod(bin, 0o755); err != nil {
		cliBinErrs[name] = err
		return "", err
	}
	cliBinPaths[name] = bin
	return bin, nil
}


func buildGoBestPracticeOnce(t *testing.T, d *session.Doctest) (string, error) {
	return buildCLIBinaryOnce(t, d, "go-best-practice", "./cmd/go-best-practice")
}

func buildPlaywrightDebugOnce(t *testing.T, d *session.Doctest) (string, error) {
	return buildCLIBinaryOnce(t, d, "playwright-debug", "./cmd/playwright-debug")
}

func buildGithubFetchOnce(t *testing.T, d *session.Doctest) (string, error) {
	return buildCLIBinaryOnce(t, d, "github-fetch", "./cmd/github-fetch")
}

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
