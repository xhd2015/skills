# Scenario

**Feature**: skills `cmd/go-best-practice` is a redirect stub to the standalone module

```
# any argv hits the same redirect (no recipes / skill / vet)
user -> go-best-practice [args...] -> stderr redirect + exit 1

# new home
redirect message -> github.com/xhd2015/go-best-practice
redirect message -> go install github.com/xhd2015/go-best-practice/cmd/go-best-practice@latest
```

## Preconditions

- Module root contains `go.mod` for `github.com/xhd2015/skills`.
- Binary is built once per `doctest test` session into a temp cache keyed by
  `d.DOCTEST_SESSION_ID` from `./cmd/go-best-practice`.
- Child process env includes `GOWORK=off` (parallel-safe; no product stdout hijack).

## Steps

1. Session-cache build `./cmd/go-best-practice`.
2. Each leaf sets `req.Args` (may be empty).
3. `Run` executes the binary and captures stdout, stderr, and exit code.

## Context

- Classic TDD RED until implementer replaces full CLI with stub.
- Prefer unlabeled leaves (small suite, ≤3 leaves; no `e2e` label needed per P4 plan).
- Shared `assertRedirect` lives in `DOCTEST.md` Go block.
- Process cwd is undetermined — resolve module root via `d.DOCTEST_ROOT`.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func moduleRoot(d *session.Doctest) (string, error) {
	dir := d.DOCTEST_ROOT
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot find go.mod from %s", d.DOCTEST_ROOT)
		}
		dir = parent
	}
}

func sessionCacheDir(d *session.Doctest) string {
	return filepath.Join(os.TempDir(), "go-best-practice-redirect-"+d.DOCTEST_SESSION_ID)
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

func buildGoBestPracticeOnce(t *testing.T, d *session.Doctest) (string, error) {
	t.Helper()
	cacheDir := sessionCacheDir(d)
	bin := filepath.Join(cacheDir, "go-best-practice")
	ready := filepath.Join(cacheDir, "binaries.ready")
	lock := filepath.Join(cacheDir, "build.lock")

	err := withFileLock(t, lock, func() error {
		if _, err := os.Stat(ready); err == nil {
			if _, err := os.Stat(bin); err == nil {
				return nil
			}
		}
		root, err := moduleRoot(d)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", bin, "./cmd/go-best-practice")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GOWORK=off")
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build go-best-practice: %w\n%s", err, out)
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

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
