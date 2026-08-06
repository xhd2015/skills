# Scenario

**Feature**: update command help documents supported flags

```
HandleUpdate --help -> usage text on stdout
```

## Preconditions

- No pre-install required.

## Steps

1. Leaves pass `-h` or `--help` in `UpdateArgs`.

## Context

- Help must not mutate filesystem.
- Product keeps subcommand `skills update` / usage name; no primary `--update` flag.

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
