# Scenario

**Feature**: `--help` redirects (does not print full old CLI help)

```
user -> go-best-practice --help -> stderr redirect + exit 1 (no full old help)
```

## Preconditions

- Binary is built from skills `./cmd/go-best-practice`.

## Steps

1. Set `req.Args = []string{"--help"}`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--help"}
	return nil
}
```
