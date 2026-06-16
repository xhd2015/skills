## Expected
- No error is returned.
- stdout contains `"Installed skill to:"` (indicating a fresh install, not an overwrite).
- stdout does **not** contain `"Update skill at"`.

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
	if strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout should not contain 'Update skill at' for fresh install:\n%s", resp.Stdout)
	}
}
```
