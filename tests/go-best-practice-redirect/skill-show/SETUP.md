# Scenario

**Feature**: `skill --show` redirects (does not print embedded SKILL.md body)

```
user -> go-best-practice skill --show -> stderr redirect + exit 1 (no skill body)
```

## Preconditions

- Binary is built from skills `./cmd/go-best-practice`.

## Steps

1. Set `req.Args = []string{"skill", "--show"}`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"skill", "--show"}
	return nil
}
```
