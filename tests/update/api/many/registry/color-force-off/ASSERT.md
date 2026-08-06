## Expected

- stdout has no ANSI escape sequences.
- polished status lines still present (`skill-alpha  up to date`, summary).

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
	if strings.Contains(resp.Stdout, "\x1b") {
		t.Fatalf("expected no ANSI with --no-color:\n%q", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "skill-alpha  up to date") {
		t.Fatalf("missing polished up-to-date line:\n%s", resp.Stdout)
	}
}
```