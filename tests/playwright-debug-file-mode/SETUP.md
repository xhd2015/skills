# Scenario

**Feature**: playwright-debug supports file mode and explicit eval flag routing

```
# CLI router classifies invocation mode
user -> playwright-debug CLI -> help | file runner | eval runner

# file mode executes .js with bootstrap.cjs
file runner -> bootstrap.cjs -> node + playwright -> user script stdout

# eval mode wraps snippet in async IIFE
eval runner -> node -e -> playwright -> script stdout
```

## Preconditions

- Module root is two levels above `d.DOCTEST_ROOT` (`go.mod` at `github.com/xhd2015/skills`).
- Shared fixtures live in `d.DOCTEST_ROOT/testdata/` (including `print-argv.js` for
  script-arg forwarding leaves and `print-help.js` for script `--help` forwarding).
- The `playwright-debug` binary is built once per `doctest test` session into a
  temp cache keyed by `DOCTEST_SESSION_ID`.
- Fast routing leaves expect failures **before** playwright install; slow leaves
  require `node`, `npm`, and network for first-time playwright cache setup.

## Steps

1. Build `playwright-debug` from `cmd/playwright-debug` via session cache.
2. Each leaf sets `req.Args` for its CLI invocation.
3. `Run` executes the binary and captures stdout, stderr, and exit code.

## Context

- Routing-error leaves use substring assertions on combined output.
- Execution leaves are labeled `slow` and match marker strings in stdout.
- Help leaves assert usage mentions `run`, eval flags, and file-alias behavior.

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

func buildPlaywrightDebugOnce(t *testing.T, d *session.Doctest) (string, error) {
	return buildCLIBinaryOnce(t, d, "playwright-debug", "./cmd/playwright-debug")
}

func fixturePath(d *session.Doctest, parts ...string) string {
	return filepath.Join(append([]string{d.DOCTEST_ROOT, "testdata"}, parts...)...)
}

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```