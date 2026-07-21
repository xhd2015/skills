## Expected

- No error.
- stdout contains `# Foo Skill Body` (same as flag-before-name).

## Side Effects

- None.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
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
