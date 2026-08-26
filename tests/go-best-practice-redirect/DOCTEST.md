# go-best-practice Redirect Stub Tests

Doc-style tests for the skills monorepo `cmd/go-best-practice` **redirect stub**.

After go-best-practice moved to https://github.com/xhd2015/go-best-practice
(module `github.com/xhd2015/go-best-practice`), the skills package path
`github.com/xhd2015/skills/cmd/go-best-practice` remains installable for old
bookmarks but must only print a move/redirect message (not serve recipes,
skill show, vet, or full help).

# DSN (Domain Specific Notion)

Participants:

- **User** invokes `go-best-practice` with any argv (including none, `--help`,
  `skill --show`, …) via the old skills module install path.
- **Redirect stub** (`cmd/go-best-practice`) always prints a move message and
  the new install command; it does **not** embed SKILL.md, topics, or vet.
- **stderr** carries the redirect (Error-style / deprecation so pipes stay clean).
- **stdout** stays empty (no skill body, no old full help).
- **exit code** is `1` for every invocation.

Behaviors:

- Any args (including empty) → same redirect UX.
- Output must mention `https://github.com/xhd2015/go-best-practice`, module
  `github.com/xhd2015/go-best-practice`, and the install line:
  `go install github.com/xhd2015/go-best-practice/cmd/go-best-practice@latest`.

## Decision Tree

```text
go-best-practice-redirect/
├── no-args/          # bare invocation
├── skill-show/       # skill --show must not print skill body
└── help-flag/        # --help must not print full old help
```

## Test Index

| Leaf | Args | Expected |
|------|------|----------|
| `no-args` | `[]` | exit 1; stderr contains module + `go install …@latest`; no product body on stdout |
| `skill-show` | `skill --show` | same redirect; no SKILL.md body |
| `help-flag` | `--help` | same redirect; no full old CLI help |

## How to Run

```sh
doctest vet ./tests/go-best-practice-redirect
doctest test -v ./tests/go-best-practice-redirect
```

## Version

0.0.2

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// Expected redirect fragments (must appear on stderr).
const (
	redirectURL     = "https://github.com/xhd2015/go-best-practice"
	redirectModule  = "github.com/xhd2015/go-best-practice"
	redirectInstall = "go install github.com/xhd2015/go-best-practice/cmd/go-best-practice@latest"
)

type Request struct {
	Args []string
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	bin, err := buildGoBestPracticeOnce(t, d)
	if err != nil {
		return nil, err
	}

	args := req.Args
	if args == nil {
		args = []string{}
	}

	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if cmd.ProcessState != nil {
		resp.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runErr != nil {
		if _, ok := runErr.(*exec.ExitError); !ok {
			return resp, fmt.Errorf("run go-best-practice: %w", runErr)
		}
	}
	return resp, nil
}

// assertRedirect checks the shared stub contract: exit 1, stderr redirect,
// empty product stdout.
func assertRedirect(t *testing.T, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Run failed: %v\nstdout:\n%s\nstderr:\n%s", err, resp.Stdout, resp.Stderr)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("exit code = %d, want 1\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, redirectURL) {
		t.Fatalf("stderr missing URL %q:\n%s", redirectURL, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, redirectModule) {
		t.Fatalf("stderr missing module %q:\n%s", redirectModule, resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, redirectInstall) {
		t.Fatalf("stderr missing install line %q:\n%s", redirectInstall, resp.Stderr)
	}
	if strings.TrimSpace(resp.Stdout) != "" {
		t.Fatalf("stdout must be empty for redirect stub (no skill body / old help):\n%s", resp.Stdout)
	}
}
```
