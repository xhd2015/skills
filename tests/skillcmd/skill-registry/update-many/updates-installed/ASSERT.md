## Expected

- No error.
- stdout indicates update or up-to-date for foo (after restore path, update message preferred).
- On-disk foo SKILL.md matches canonical foo content after update.
- bar reports not installed.

## Side Effects

- foo SKILL.md content restored to registry content.

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
	if !strings.Contains(resp.Stdout, "Update skill at") && !strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("expected update or up-to-date line for foo:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "skill not installed: bar") {
		t.Fatalf("expected not-installed for bar:\n%s", resp.Stdout)
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
