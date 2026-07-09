## Expected

- No error.
- stdout contains `foo` and `bar`.
- stdout contains description text for at least one skill when present.

## Side Effects

- None.

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
	if !strings.Contains(resp.Stdout, "foo") {
		t.Fatalf("stdout missing foo:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "bar") {
		t.Fatalf("stdout missing bar:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "foo skill description") && !strings.Contains(resp.Stdout, "bar skill description") {
		t.Fatalf("stdout should include at least one description:\n%s", resp.Stdout)
	}
}
```
