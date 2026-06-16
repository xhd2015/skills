## Expected
- No error is returned.
- stdout contains `"Update skill at"` (because the directory existed and was overwritten).
- stdout does **not** contain `"Aborted."` (confirmation is skipped when `--force` is present).

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
	if !strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout missing 'Update skill at':\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Aborted.") {
		t.Fatalf("stdout should not contain 'Aborted.' when --force is set:\n%s", resp.Stdout)
	}
}
```
