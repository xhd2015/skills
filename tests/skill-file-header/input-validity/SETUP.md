# Scenario

**Feature**: header extraction succeeds only when SKILL.md opens with YAML frontmatter

```
# valid content exposes inner YAML
SKILL.md content -> GetHeader -> inner YAML

# missing delimiter prefix is rejected
SKILL.md content -> GetHeader -> error
```

## Preconditions

- Fixture files live alongside each leaf as `skill_content.md`.

## Steps

1. Leaves load `skill_content.md` into `req.Content` before `Run`.

## Context

- `ParseHeader` is exercised only on the valid-header leaf where `GetHeader` succeeds.

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
	if req == nil {
		t.Fatal("request must be initialized")
	}
	return nil
}
```