## Expected

- stdout contains ANSI escape sequences (`\x1b`).
- stdout still contains plain skill names and status words when escapes are ignored conceptually (name + "up to date" / "not installed").
- green SGR (`\x1b[32m`) is not required for up-to-date-only runs; gray (`\x1b[90m`) is expected for status/summary.

## Errors

- None.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "\x1b[") {
		t.Fatalf("expected ANSI escapes with --color:\n%q", resp.Stdout)
	}
	// Gray for up to date / not installed / summary
	if !strings.Contains(resp.Stdout, "\x1b[90m") {
		t.Fatalf("expected gray SGR for status/summary:\n%q", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "skill-alpha") {
		t.Fatalf("missing skill-alpha:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "up to date") {
		t.Fatalf("missing up to date status:\n%s", resp.Stdout)
	}
}
```
