# Scenario

**Feature**: repo skill CLIs support `skill --install` with dry-run and global targets

```
# user invokes skill --install on a CLI binary
user -> skill CLI (skill --install) -> install/skillcmd HandleInstall -> stdout / filesystem

# global scope resolves under HOME
user -> skill CLI (skill --install --global) -> ~/.agents/skills/<name>

# go-best-practice ships nested skill-cli/SKILL.md extras (Shape 3)
HandleInstall -> SKILL.md + skill-cli/SKILL.md (+ other nested paths)
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
	return filepath.Join(os.TempDir(), "skill-install-cli-"+DOCTEST_SESSION_ID)
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

func buildCLIBinaryOnce(t *testing.T, name, pkg string) (string, error) {
	t.Helper()
	cacheDir := sessionCacheDir()
	bin := filepath.Join(cacheDir, name)
	ready := filepath.Join(cacheDir, name+".ready")
	lock := filepath.Join(cacheDir, name+".lock")

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
		cmd := exec.Command("go", "build", "-o", bin, pkg)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "GOWORK=off")
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build %s: %w\n%s", name, err, out)
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

func buildGoBestPracticeOnce(t *testing.T) (string, error) {
	return buildCLIBinaryOnce(t, "go-best-practice", "./cmd/go-best-practice")
}

func buildPlaywrightDebugOnce(t *testing.T) (string, error) {
	return buildCLIBinaryOnce(t, "playwright-debug", "./cmd/playwright-debug")
}

func buildGithubFetchOnce(t *testing.T) (string, error) {
	return buildCLIBinaryOnce(t, "github-fetch", "./cmd/github-fetch")
}

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```
