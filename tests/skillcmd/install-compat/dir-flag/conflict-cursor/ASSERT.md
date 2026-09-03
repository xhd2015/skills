## Expected

- Error that `--dir` cannot combine with preset target flags.
- No install under `.cursor` or collection.

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
		t.Fatal("expected error when --dir and --cursor both set")
	}
	if !strings.Contains(resp.Error, "--dir") {
		t.Fatalf("error=%q", resp.Error)
	}
	if _, e := os.Stat(filepath.Join(resp.WorkDir, ".cursor", "skills", "demo-skill", "SKILL.md")); e == nil {
		t.Fatal("must not install to --cursor on conflict")
	}
}
```
