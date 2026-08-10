## Expected
- Exit code 0.
- stdout shows auto-detected repo and `State: open`.
- stdout lists issues #15 and #14 but NOT PR #42.
- stdout contains `Page 1 (2 items)`.

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
		"#15", "@alice", "Login page broken",
		"#14", "@bob", "Add dark mode",
		"Page 1 (2 items)",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "#42") || strings.Contains(stdout, "Fix login redirect") {
		t.Fatalf("stdout should not contain PR masquerading as issue:\n%s", stdout)
	}
}
```