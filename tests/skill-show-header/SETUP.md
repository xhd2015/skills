# Scenario

**Feature**: skill CLI binaries support `skill --show --header`

```
# user invokes skill --show on a CLI binary
user -> skill CLI -> embedded SKILL.md content

# header-only mode trims body
user -> skill CLI (skill --show --header) -> YAML frontmatter only
```

## Preconditions

- The `go-best-practice` skill CLI is built from this module.
- The embedded SKILL.md body contains the marker `# Go Best Practice Skill`.

## Steps

1. Build `go-best-practice` into a temporary binary path on `req.Binary`.
2. Leaves set header mode flags or explicit `req.Args`.

## Context

- Tests use the real `go-best-practice` binary to exercise skill flag parsing
  after migration to skillcmd.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Binary != "" {
		return nil
	}
	repoRoot := filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	req.Binary = filepath.Join(t.TempDir(), "go-best-practice")
	cmd := exec.Command("go", "build", "-o", req.Binary, "./cmd/go-best-practice")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build go-best-practice: %w\n%s", err, out)
	}
	if err := os.Chmod(req.Binary, 0o755); err != nil {
		return err
	}
	return nil
}
```
