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
	"testing"
)

func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return string(data)
}

func Setup(t *testing.T, req *Request) error {
	if req.Content != "" {
		t.Fatalf("req.Content = %q before leaf setup, want empty", req.Content)
	}
	return nil
}
```