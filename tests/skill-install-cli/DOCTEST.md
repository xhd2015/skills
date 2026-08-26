# Skill Install CLI Tests

Doc-style tests for `skill --install` on repo skill CLIs that still ship
embedded skill content (`playwright-debug`, `github-fetch`). Covers dry-run
output, global target resolution, real install side effects, and routing
rules. `go-best-practice` install/show coverage lives in the standalone
module https://github.com/xhd2015/go-best-practice (`tests/skill-install-cli`).

# DSN (Domain Specific Notion)

Participants:

- User invokes a skill CLI binary with `skill --install`, top-level `install`, or
  `skill --show` (actions are **flags**, not word subcommands `show`/`install`).
- CLI router (`handle`) dispatches `skill` into skillcmd (or equivalent) which
  classifies `--show` / `--install` / `--list` and remaining path/name/install flags.
- Install resolves target directories (local `.agents/skills/<name>`, global
  `~/.agents/skills/<name>`, or explicit `<dir>`) and copies SKILL.md plus any
  nested extra files.
- Dry-run mode prints `[dry-run]` lines to stdout without writing files.
- Legacy word forms (`skill show`, `skill install`) and bare paths under `skill`
  without an action flag are rejected.

## Decision Tree

```text
skill-install-cli/
├── playwright-debug/
│   └── global-dry-run/
└── github-fetch/
    ├── global-dry-run/
    └── standalone-install-rejected/
```

## Test Index

| Leaf | Description |
|------|-------------|
| `playwright-debug/global-dry-run` | Global dry-run mentions `playwright-debug` under `HOME` |
| `github-fetch/global-dry-run` | Global dry-run mentions `github-fetch` under `HOME` |
| `github-fetch/standalone-install-rejected` | Top-level `install` is unknown (must use `skill --install`) |

## How to Run

```sh
doctest vet ./tests/skill-install-cli
doctest test -v ./tests/skill-install-cli
```

## Version

0.0.3

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

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

	resp := &Response{}
	env := append([]string{}, os.Environ()...)
	env = append(env, "GOWORK=off")

	cmd := exec.Command(req.Binary, req.Args...)
	if req.UseWorkDir {
		workDir := t.TempDir()
		resp.WorkDir = workDir
		cmd.Dir = workDir
	}
	if req.UseGlobalHome {
		homeDir := t.TempDir()
		resp.HomeDir = homeDir
		env = append(env, "HOME="+homeDir)
	}
	cmd.Env = env

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
