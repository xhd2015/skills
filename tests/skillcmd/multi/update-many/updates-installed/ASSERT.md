## Expected

- No error.
- stdout reports polished `foo  updated` (with file lines) after restoring drifted content.
- On-disk foo SKILL.md matches canonical foo content after update.
- bar reports `not installed` with polished summary.

## Side Effects

- foo SKILL.md content restored to registry content.

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
	if !strings.Contains(resp.Stdout, "foo  updated") {
		t.Fatalf("expected polished updated status for foo:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "bar  not installed") {
		t.Fatalf("expected not-installed for bar:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "1 updated · 0 up to date · 1 not installed") {
		t.Fatalf("missing summary line:\n%s", resp.Stdout)
	}
	p := filepath.Join(resp.WorkDir, ".agents", "skills", "foo", "SKILL.md")
	data, readErr := os.ReadFile(p)
	if readErr != nil {
		t.Fatalf("read foo SKILL.md: %v", readErr)
	}
	if !strings.Contains(string(data), "# Foo Skill Body") {
		t.Fatalf("foo SKILL.md not restored to canonical content:\n%s", data)
	}
}
```
