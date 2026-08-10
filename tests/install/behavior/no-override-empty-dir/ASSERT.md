## Expected
- No error is returned.
- stdout contains `"Update skill at"` (directory already existed, even though empty).
- stdout does **not** contain `"Aborted."`.
- Install still proceeds without confirmation for an empty dir under `--no-override`.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	// Dir existed before install → Update header (not Installed skill to).
	if !strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout missing 'Update skill at' for pre-existing empty dir:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Aborted.") {
		t.Fatalf("stdout should not contain 'Aborted.' for empty directory:\n%s", resp.Stdout)
	}
}
```
