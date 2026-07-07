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

- Module root is two levels above `DOCTEST_ROOT` (`go.mod` at `github.com/xhd2015/skills`).
- Shared fixtures live in `DOCTEST_ROOT/testdata/`.
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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

func moduleRoot() (string, error) {
	dir := DOCTEST_ROOT
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot find go.mod from %s", DOCTEST_ROOT)
		}
		dir = parent
	}
}

func sessionCacheDir() string {
	return filepath.Join(os.TempDir(), "playwright-debug-file-mode-"+DOCTEST_SESSION_ID)
}

func withFileLock(t *testing.T, lockPath string, fn func() error) error {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return fn()
}

func buildPlaywrightDebugOnce(t *testing.T) (string, error) {
	t.Helper()
	cacheDir := sessionCacheDir()
	bin := filepath.Join(cacheDir, "playwright-debug")
	ready := filepath.Join(cacheDir, "binaries.ready")
	lock := filepath.Join(cacheDir, "build.lock")

	err := withFileLock(t, lock, func() error {
		if _, err := os.Stat(ready); err == nil {
			if _, err := os.Stat(bin); err == nil {
				return nil
			}
		}
		root, err := moduleRoot()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", bin, "./cmd/playwright-debug")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GOWORK=off")
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build playwright-debug: %w\n%s", err, out)
		}
		if err := os.Chmod(bin, 0o755); err != nil {
			return err
		}
		return os.WriteFile(ready, []byte("ok"), 0o644)
	})
	if err != nil {
		return "", err
	}
	return bin, nil
}

func fixturePath(parts ...string) string {
	return filepath.Join(append([]string{DOCTEST_ROOT, "testdata"}, parts...)...)
}

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```