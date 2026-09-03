## Expected

- No error.
- Dry-run targets `vendor/skills/demo-skill`.
- No "unexpected arguments" from dual positionals.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s\nstdout=%s", resp.Error, resp.Stdout)
	}
	want := filepath.Join("vendor", "skills", "demo-skill")
	if !strings.Contains(resp.Stdout, want) {
		t.Fatalf("stdout missing nested skill root %q:\n%s", want, resp.Stdout)
	}
	if strings.Contains(strings.ToLower(resp.Error), "unexpected arguments") {
		t.Fatalf("topic must be peeled:\n%s", resp.Error)
	}
}
```
