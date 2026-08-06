## Expected

- Nested topic body and frontmatter name are printed from `a/b/TOPIC.md`.
- Root body is not printed.

## Side Effects

- None (show is read-only).

## Errors

- None.

```go
import (
	"github.com/xhd2015/doctest/session"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "# Nested A/B Body") {
		t.Fatalf("stdout missing nested body:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "name: demo-skill/a/b") {
		t.Fatalf("stdout missing nested name:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "# Demo Skill Body") {
		t.Fatalf("must not print root body for nested path:\n%s", resp.Stdout)
	}
}
```
