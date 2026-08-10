## Expected
- No error is returned.
- stdout contains `"Update skill at"` (indicating an overwrite, not a fresh install).
- stdout does **not** contain `"Installed skill to:"`.

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
	if !strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout missing 'Update skill at':\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Installed skill to:") {
		t.Fatalf("stdout should not contain 'Installed skill to:' for overwrite:\n%s", resp.Stdout)
	}
}
```
