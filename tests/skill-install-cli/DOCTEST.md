# Skill Install CLI Tests

Doc-style tests for `skill --install` on the three repo skill CLIs (`go-best-practice`,
`playwright-debug`, `github-fetch`). Covers dry-run output, global target resolution,
real install side effects (nested `path/TOPIC.md` extras), backward-compatible
top-level `install`, regression for `skill --show`, nested topic show orders, and
error paths for missing/legacy action forms.

# DSN (Domain Specific Notion)

Participants:

- User invokes a skill CLI binary with `skill --install`, top-level `install`, or
  `skill --show` (actions are **flags**, not word subcommands `show`/`install`).
- CLI router (`handle`) dispatches `skill` into skillcmd (or equivalent) which
  classifies `--show` / `--install` / `--list` and remaining path/name/install flags.
- Install resolves target directories (local `.agents/skills/<name>`, global
  `~/.agents/skills/<name>`, or explicit `<dir>`) and copies SKILL.md plus any
  nested extra files (e.g. `cli/skill-cli/TOPIC.md` for go-best-practice Shape 3).
- Dry-run mode prints `[dry-run]` lines to stdout without writing files.
- Legacy word forms (`skill show`, `skill install`) and bare paths under `skill`
  without an action flag are rejected.

## Decision Tree

```text
skill-install-cli/
├── go-best-practice/
│   ├── skill-install/
│   │   ├── dry-run-default/
│   │   ├── global-dry-run/
│   │   └── includes-topics/          # nested cli/skill-cli/TOPIC.md (not topics/*.md)
│   ├── skill-show-still-works/
│   ├── skill-show-nested/
│   │   ├── flag-before-path/
│   │   └── path-before-flag/
│   ├── bare-path-no-action/
│   ├── legacy-skill-show/
│   ├── top-level-install-alias/
│   └── unknown-subcommand/           # skill with no action flags
├── playwright-debug/
│   └── global-dry-run/
└── github-fetch/
    ├── global-dry-run/
    └── standalone-install-rejected/
```

## Test Index

| Leaf | Description |
|------|-------------|
| `go-best-practice/skill-install/dry-run-default` | `skill --install --dry-run` mentions `.agents/skills/go-best-practice` |
| `go-best-practice/skill-install/global-dry-run` | `skill --install --global --dry-run` resolves under `HOME` |
| `go-best-practice/skill-install/includes-topics` | Real install copies `cli/skill-cli/TOPIC.md` with SKILL.md |
| `go-best-practice/skill-show-still-works` | `skill --show` still prints skill name (regression) |
| `go-best-practice/skill-show-nested/flag-before-path` | `skill --show cli/skill-cli` prints nested topic; name has `go-best-practice/cli/skill-cli` |
| `go-best-practice/skill-show-nested/path-before-flag` | `skill cli/skill-cli --show` same nested content |
| `go-best-practice/bare-path-no-action` | `skill cli/skill-cli` without action → error |
| `go-best-practice/legacy-skill-show` | legacy `skill show` → error |
| `go-best-practice/top-level-install-alias` | Top-level `install --dry-run` matches `skill --install --dry-run` |
| `go-best-practice/unknown-subcommand` | `skill` alone (no action flags) errors with expected action hints |
| `playwright-debug/global-dry-run` | Global dry-run mentions `playwright-debug` under `HOME` |
| `github-fetch/global-dry-run` | Global dry-run mentions `github-fetch` under `HOME` |
| `github-fetch/standalone-install-rejected` | Top-level `install` is unknown (must use `skill --install`) |

## How to Run

```sh
doctest vet ./tests/skill-install-cli
doctest test -v ./tests/skill-install-cli
```

## Version

0.0.2

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// withProcessLock serializes process-global chdir across workspace trees.
func withProcessLock(fn func() error) error {
	lockPath := filepath.Join(os.TempDir(), "skills-doctest-process.lock")
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

type Request struct {
	Binary        string
	SkillName     string
	Args          []string
	UseGlobalHome bool
	UseWorkDir    bool
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	WorkDir  string
	HomeDir  string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	if req.Binary == "" {
		return nil, fmt.Errorf("CLI binary path is required")
	}
	if len(req.Args) == 0 {
		return nil, fmt.Errorf("CLI args are required")
	}

	var resp *Response
	err := withProcessLock(func() error {
		var runErr error
		resp, runErr = runLocked(t, req)
		return runErr
	})
	return resp, err
}

func runLocked(t *testing.T, req *Request) (*Response, error) {
	resp := &Response{}

	if req.UseWorkDir {
		workDir := t.TempDir()
		resp.WorkDir = workDir
		prevWD, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		if err := os.Chdir(workDir); err != nil {
			return nil, err
		}
		defer func() {
			_ = os.Chdir(prevWD)
		}()
	}

	cmd := exec.Command(req.Binary, req.Args...)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if req.UseGlobalHome {
		homeDir := t.TempDir()
		resp.HomeDir = homeDir
		// Parallel-safe: set HOME on child env only (not t.Setenv).
		cmd.Env = append(cmd.Env, "HOME="+homeDir)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	resp.Stdout = stdout.String()
	resp.Stderr = stderr.String()
	if cmd.ProcessState != nil {
		resp.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runErr != nil {
		if _, ok := runErr.(*exec.ExitError); !ok {
			return resp, fmt.Errorf("run %s: %w", filepath.Base(req.Binary), runErr)
		}
	}
	return resp, nil
}
```
