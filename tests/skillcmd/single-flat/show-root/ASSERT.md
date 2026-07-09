## Expected

- No error.
- See scenario-specific stdout checks in Assert.

## Side Effects

- install-dry-run must not create directories.

## Errors

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
	if !strings.Contains(resp.Stdout, "name: demo-skill") {
		t.Fatalf("stdout missing name header:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "# Demo Skill Body") {
		t.Fatalf("stdout missing body marker:\n%s", resp.Stdout)
	}
}
```
