# Skill Show Header CLI Tests

Doc-style tests for `skill show --header` on skill CLI binaries. The flag must
print only YAML frontmatter (including delimiters) and omit the Markdown body.

# DSN (Domain Specific Notion)

Participants:

- A skill CLI binary embeds the full SKILL.md template and exposes `skill show`
  to print the entire file.
- `skill show --header` prints only the YAML frontmatter block with opening and
  closing `---` delimiter lines.
- Without `--header`, `skill show` keeps the existing behavior of printing the
  full embedded SKILL.md content including the body.

## Decision Tree

```text
skill-show-header
└── show-mode
    ├── with-header-flag
    └── without-header-flag
```

## Test Index

- `show-mode/with-header-flag`: `skill show --header` prints delimiters and `name:` but not the body-only marker.
- `show-mode/without-header-flag`: `skill show` prints both header fields and the body marker.

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
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type Request struct {
	Binary     string
	HeaderOnly bool
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Binary == "" {
		return nil, fmt.Errorf("skill CLI binary path is required")
	}
	args := []string{"skill", "show"}
	if req.HeaderOnly {
		args = append(args, "--header")
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