## Expected

- FormatErr empty.
- Formatted starts with `---` and contains `name: demo-skill`.
- Formatted does not contain `# Demo Skill Body`.

## Side Effects

- None.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.FormatErr != "" {
		t.Fatalf("FormatHeader error: %s", resp.FormatErr)
	}
	if !strings.HasPrefix(resp.Formatted, "---") {
		t.Fatalf("Formatted must start with ---: %q", resp.Formatted)
	}
	if !strings.Contains(resp.Formatted, "name: demo-skill") {
		t.Fatalf("Formatted missing name:\n%s", resp.Formatted)
	}
	if strings.Contains(resp.Formatted, "# Demo Skill Body") {
		t.Fatalf("Formatted must omit body:\n%s", resp.Formatted)
	}
}
```
