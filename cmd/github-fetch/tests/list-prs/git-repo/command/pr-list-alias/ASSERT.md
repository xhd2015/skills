## Expected
- Exit code 0.
- stdout matches the `open-default` case: auto-detected repo, open state, both PRs, and `Page 1 (2 items)`.

## Exit Code
- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	stdout := resp.Stdout
	for _, want := range []string{
		"testowner/testrepo",
		"https://github.com/testowner/testrepo",
		"State: open",
		"#42", "@alice", "Fix login redirect",
		"#41", "@bob", "Add pagination to API client",
		"Page 1 (2 items)",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
}
```