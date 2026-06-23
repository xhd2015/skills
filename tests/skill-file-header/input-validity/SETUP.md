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
	if req == nil {
		t.Fatal("request must be initialized")
	}
	return nil
}
```