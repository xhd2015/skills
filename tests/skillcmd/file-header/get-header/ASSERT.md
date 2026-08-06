## Expected

- HeaderErr empty.
- HeaderText is `name: git-fetch\ndescription: clone repos` (no `---`).

## Side Effects

- None.

```go
import (
	"github.com/xhd2015/doctest/session"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.HeaderErr != "" {
		t.Fatalf("GetHeader error: %s", resp.HeaderErr)
	}
	want := "name: git-fetch\ndescription: clone repos"
	if strings.TrimSpace(resp.HeaderText) != want {
		t.Fatalf("HeaderText = %q, want %q", resp.HeaderText, want)
	}
	if strings.Contains(resp.HeaderText, "---") {
		t.Fatalf("HeaderText must not include delimiters: %q", resp.HeaderText)
	}
}
```
