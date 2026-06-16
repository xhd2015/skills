## Expected
- No error is returned.
- stdout contains `"Installed skill to:"` (the empty directory is treated as a fresh install).
- stdout does **not** contain `"Aborted."`.

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
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "Installed skill to:") {
		t.Fatalf("stdout missing 'Installed skill to:':\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Aborted.") {
		t.Fatalf("stdout should not contain 'Aborted.' for empty directory:\n%s", resp.Stdout)
	}
}
```
