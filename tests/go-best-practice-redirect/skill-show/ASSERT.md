## Expected

- Exit code `1`.
- stderr contains module + new `go install` path (shared redirect).
- stdout is empty — no SKILL.md body (e.g. no "Go Best Practice Skill" index).

## Side Effects

- None beyond process I/O.

## Errors

- No error from `Run`.

## Exit Code

- 1

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertRedirect(t, resp, err)
	// Extra guard: old full CLI skill --show printed this heading on stdout.
	if strings.Contains(resp.Stdout, "Go Best Practice Skill") {
		t.Fatalf("stdout must not contain skill body:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "name: go-best-practice") {
		t.Fatalf("stdout must not contain SKILL.md frontmatter:\n%s", resp.Stdout)
	}
}
```
