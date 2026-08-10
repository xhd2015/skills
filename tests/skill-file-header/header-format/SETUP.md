# Scenario

**Feature**: YAML header parsing supports folded block scalar values

```
# folded description uses YAML >- syntax
header YAML -> ParseHeader -> single normalized description value
```

## Preconditions

- `GetHeader` succeeds on fixture content with a folded `description` field.

## Steps

1. Leaves load `skill_content.md` into `req.Content`.

## Context

- Folded scalars must collapse embedded newlines into spaces for lookup via `Entries.Get`.

```go
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func readFixture(t *testing.T, d *session.Doctest, name string) string {
	t.Helper()
	// DOCTEST_CASE is absolute path to the leaf directory.
	path := filepath.Join(d.DOCTEST_CASE, name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", path, err)
	}
	return string(data)
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Content != "" {
		t.Fatalf("req.Content = %q before leaf setup, want empty", req.Content)
	}
	return nil
}
```