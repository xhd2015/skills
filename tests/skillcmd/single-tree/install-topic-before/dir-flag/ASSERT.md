## Expected

- No error.
- Dry-run targets `vendor/skills/demo-skill` (topic peeled; collection nested).
- Does not treat `skill-cli` as the install directory.

```go
import (
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
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s\nstdout=%s", resp.Error, resp.Stdout)
	}
	want := filepath.Join("vendor", "skills", "demo-skill")
	if !strings.Contains(resp.Stdout, want) {
		t.Fatalf("stdout missing nested skill root %q:\n%s", want, resp.Stdout)
	}
	if strings.Contains(resp.Stdout, filepath.Join("skill-cli", "SKILL.md")) &&
		!strings.Contains(resp.Stdout, want) {
		t.Fatalf("treated topic as install dir:\n%s", resp.Stdout)
	}
	if strings.Contains(strings.ToLower(resp.Error+resp.Stdout), "not both") {
		t.Fatalf("topic must be peeled so --dir is alone:\n%s%s", resp.Error, resp.Stdout)
	}
}
```
