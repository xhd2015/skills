# Skill Install CLI Tests

Doc-style tests for `skill install` on the three repo skill CLIs (`go-best-practice`,
`playwright-debug`, `github-fetch`). Covers dry-run output, global target resolution,
real install side effects, backward-compatible top-level `install`, regression for
`skill show`, and error paths for unknown sub-commands.

# DSN (Domain Specific Notion)

Participants:

- User invokes a skill CLI binary with `skill install`, top-level `install`, or
  `skill show`.
- CLI router (`handle`) dispatches `skill` to `handleSkill`, which wires `install`
  for `skill install` or prints embedded SKILL.md for `skill show`.
- `install.HandleInstall` resolves target directories (local `.agents/skills/<name>`,
  global `~/.agents/skills/<name>`, or explicit `<dir>`) and copies SKILL.md plus
  any embedded extra files (topics for `go-best-practice`).
- Dry-run mode prints `[dry-run]` lines to stdout without writing files.

## Decision Tree

```text
skill-install-cli/
├── go-best-practice/
│   ├── skill-install/
│   │   ├── dry-run-default/
│   │   ├── global-dry-run/
│   │   └── includes-topics/
│   ├── skill-show-still-works/
│   ├── top-level-install-alias/
│   └── unknown-subcommand/
├── playwright-debug/
│   └── global-dry-run/
└── github-fetch/
    ├── global-dry-run/
    └── standalone-install-rejected/
```

## Test Index

| Leaf | Description |
|------|-------------|
| `go-best-practice/skill-install/dry-run-default` | `skill install --dry-run` mentions `.agents/skills/go-best-practice` |
| `go-best-practice/skill-install/global-dry-run` | `skill install --global --dry-run` resolves under `HOME` |
| `go-best-practice/skill-install/includes-topics` | Real install copies `topics/skill-cli.md` with SKILL.md |
| `go-best-practice/skill-show-still-works` | `skill show` still prints skill name (regression) |
| `go-best-practice/top-level-install-alias` | Top-level `install --dry-run` matches `skill install --dry-run` |
| `go-best-practice/unknown-subcommand` | `skill bogus` errors with known sub-command hints |
| `playwright-debug/global-dry-run` | Global dry-run mentions `playwright-debug` under `HOME` |
| `github-fetch/global-dry-run` | Global dry-run mentions `github-fetch` under `HOME` |
| `github-fetch/standalone-install-rejected` | Top-level `install` is unknown (must use `skill install`) |

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
	"testing"
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

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Binary == "" {
		return nil, fmt.Errorf("CLI binary path is required")
	}
	if len(req.Args) == 0 {
		return nil, fmt.Errorf("CLI args are required")
	}

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

	if req.UseGlobalHome {
		homeDir := t.TempDir()
		resp.HomeDir = homeDir
		t.Setenv("HOME", homeDir)
	}

	cmd := exec.Command(req.Binary, req.Args...)
	cmd.Env = append(os.Environ(), "GOWORK=off")
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