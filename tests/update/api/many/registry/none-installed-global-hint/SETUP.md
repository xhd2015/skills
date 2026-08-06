# Scenario

**Feature**: `--global` batch update with zero installs reports each skill

```
HOME temp dir, no installs -> HandleUpdateMany(..., --global)
  -> not-installed line per skill + summary
```

## Preconditions

- `req.UseGlobalHome` isolates `HOME` with no skill trees.

## Steps

1. Enable global home temp dir.
2. Run batch update with `["--global"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.UseGlobalHome = true
	req.PreInstalls = nil
	req.UpdateArgs = []string{"--global"}
	return nil
}
```
