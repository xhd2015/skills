## Expected

- `GetHeader` succeeds without error.
- `GetHeader` returns `name: git-fetch\ndescription: clone repos` with no `---` delimiters.
- `ParseHeader` succeeds without error.
- Parsed entries include `name=git-fetch` and `description=clone repos`.
- `Entries.Get("name")` returns `git-fetch`.

## Side Effects

- None; pure function calls only.

## Errors

- `GetHeader` and `ParseHeader` must not return errors.

## Exit Code

- Not applicable.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.HeaderErr != nil {
		t.Fatalf("GetHeader error: %v", resp.HeaderErr)
	}
	if resp.ParseErr != nil {
		t.Fatalf("ParseHeader error: %v", resp.ParseErr)
	}
	wantHeader := "name: git-fetch\ndescription: clone repos"
	if strings.TrimSpace(resp.Header) != wantHeader {
		t.Fatalf("GetHeader = %q, want %q", resp.Header, wantHeader)
	}
	if strings.Contains(resp.Header, "---") {
		t.Fatalf("GetHeader must not include delimiters: %q", resp.Header)
	}
	if resp.Entries.Get("name") != "git-fetch" {
		t.Fatalf(`Entries.Get("name") = %q, want git-fetch`, resp.Entries.Get("name"))
	}
	if resp.Entries.Get("description") != "clone repos" {
		t.Fatalf(`Entries.Get("description") = %q, want "clone repos"`, resp.Entries.Get("description"))
	}
}
```