## Expected

- No hard error from Run.
- stdout contains `skill not installed: foo` and `skill not installed: bar`.

## Side Effects

- No skill directories created.

```go
import (
	"os"
	"path/filepath"
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
	if !strings.Contains(resp.Stdout, "skill not installed: foo") {
		t.Fatalf("missing not-installed for foo:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "skill not installed: bar") {
		t.Fatalf("missing not-installed for bar:\n%s", resp.Stdout)
	}
	for _, name := range []string{"foo", "bar"} {
		p := filepath.Join(resp.WorkDir, ".agents", "skills", name)
		if _, statErr := os.Stat(p); !os.IsNotExist(statErr) {
			t.Fatalf("must not create install dir %s", p)
		}
	}
}
```
