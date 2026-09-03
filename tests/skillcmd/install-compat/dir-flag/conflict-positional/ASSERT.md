## Expected

- Error mentioning `--dir` / `<dir>` / not both.
- No skill files written under workdir.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error == "" {
		t.Fatal("expected error when --dir and positional both set")
	}
	lower := strings.ToLower(resp.Error)
	if !strings.Contains(lower, "--dir") || !strings.Contains(lower, "not both") {
		t.Fatalf("error=%q", resp.Error)
	}
	if _, e := os.Stat(filepath.Join(resp.WorkDir, "vendor", "skills", "demo-skill", "SKILL.md")); e == nil {
		t.Fatal("must not install on conflict")
	}
}
```
