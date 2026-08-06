## Expected

- No error.
- stdout contains `# Foo Skill Body` (same as flag-before-name).

## Side Effects

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
	if !strings.Contains(resp.Stdout, "# Foo Skill Body") {
		t.Fatalf("stdout missing foo body:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "name: foo") {
		t.Fatalf("stdout missing foo name:\n%s", resp.Stdout)
	}
}
```
