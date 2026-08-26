# Scenario

**Feature**: bare `go-best-practice` (no args) prints redirect and exits 1

```
user -> go-best-practice -> stderr redirect + exit 1
```

## Preconditions

- Binary is built from skills `./cmd/go-best-practice`.

## Steps

1. Set `req.Args = []string{}`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	return nil
}
```
