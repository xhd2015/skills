## Expected
- Exit code 0.
- stdout contains `otherowner/otherrepo` and `https://github.com/otherowner/otherrepo`.
- stdout contains `#10`, `@carol`, and `Explicit repo PR`.
- stdout shows `State: open`.

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
		"otherowner/otherrepo",
		"https://github.com/otherowner/otherrepo",
		"State: open",
		"#10",
		"@carol",
		"Explicit repo PR",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
}
```