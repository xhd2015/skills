# Scenario

**Feature**: update when no target has `SKILL.md`

```
HandleUpdate -> each resolved dir missing SKILL.md -> silent skip
```

## Preconditions

- Leaves do not run `PreInstalls`.

## Steps

1. Configure `SingleOpts` and call update.

## Context

- Sibling leaves cover the zero-install case only.
- Unlike batch, single stays silent (no `not installed` line).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.PreInstalls = nil
	return nil
}
```
