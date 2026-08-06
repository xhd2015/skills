## Expected

- No hard error from Run.
- stdout reports polished `foo  not installed` and `bar  not installed` plus a summary.

## Side Effects

- No skill directories created.

```go
import (
	"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "foo  not installed") {
		t.Fatalf("missing not-installed for foo:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "bar  not installed") {
		t.Fatalf("missing not-installed for bar:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "0 updated · 0 up to date · 2 not installed") {
		t.Fatalf("missing summary line:\n%s", resp.Stdout)
	}
	for _, name := range []string{"foo", "bar"} {
		p := filepath.Join(resp.WorkDir, ".agents", "skills", name)
		if _, statErr := os.Stat(p); !os.IsNotExist(statErr) {
			t.Fatalf("must not create install dir %s", p)
		}
	}
}
```
