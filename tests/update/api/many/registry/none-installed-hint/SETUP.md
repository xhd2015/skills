# Scenario

**Feature**: batch update with zero installed targets reports each skill

```
no pre-install -> HandleUpdateMany(registry)
  -> skill-alpha  not installed
  -> skill-beta  not installed
  -> summary with 2 not installed
```

## Preconditions

- No registry skill has `SKILL.md` at default `.agents/skills/<name>` paths.

## Steps

1. Run `HandleUpdateMany` with default target resolution (no `--global`).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.PreInstalls = nil
	req.UpdateArgs = nil
	return nil
}
```
