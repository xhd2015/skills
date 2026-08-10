# Skill Show Header CLI Tests

Doc-style tests for `skill --show --header` on skill CLI binaries. The flag must
print only YAML frontmatter (including delimiters) and omit the Markdown body.

# DSN (Domain Specific Notion)

Participants:

- A skill CLI binary embeds the full SKILL.md template and exposes `skill --show`
  to print the entire file (actions are flags, not word subcommands).
- `skill --show --header` prints only the YAML frontmatter block with opening and
  closing `---` delimiter lines.
- Without `--header`, `skill --show` keeps printing the full embedded SKILL.md
  content including the body.
- Both flag orders for header mode are valid: `--show --header` and
  `--header --show`.

## Decision Tree

```text
skill-show-header
└── show-mode
    ├── with-header-flag
    ├── without-header-flag
    └── header-before-show
```

## Test Index

- `show-mode/with-header-flag`: `skill --show --header` prints delimiters and `name:` but not the body-only marker.
- `show-mode/without-header-flag`: `skill --show` prints both header fields and the body marker.
- `show-mode/header-before-show`: `skill --header --show` matches header-only output (both orders).

## How to Run

```sh
doctest vet ./tests/skill-show-header
doctest test -v ./tests/skill-show-header
```

## Version

0.0.2

```go
import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	Binary string
	// Args is the full argv after the binary name (e.g. skill --show --header).
	// When empty, Run defaults from HeaderOnly / HeaderBeforeShow.
	Args []string

	HeaderOnly       bool
	HeaderBeforeShow bool
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	if req.Binary == "" {
		return nil, fmt.Errorf("skill CLI binary path is required")
	}
	args := req.Args
	if len(args) == 0 {
		if req.HeaderOnly && req.HeaderBeforeShow {
			args = []string{"skill", "--header", "--show"}
		} else if req.HeaderOnly {
			args = []string{"skill", "--show", "--header"}
		} else {
			args = []string{"skill", "--show"}
		}
	}
	cmd := exec.Command(req.Binary, args...)
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
	return resp, runErr
}
```
