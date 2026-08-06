# Scenario

**Feature**: `HandleUpdate` for one embedded skill

```
test harness -> HandleUpdate(InstallOptions, updateArgs)
  -> inventory apply on installed targets only
  -> polished status line (+ file lines when changed)
```

## Preconditions

- `req.UseMany` is false for all leaves under this branch.

## Steps

1. Set `req.UseMany = false`.
2. Leaves configure `req.SingleOpts` and `req.UpdateArgs`.

## Context

- Default target when no location flag is passed matches install:
  `.agents/skills/<SkillDirName>`.
- Single-skill uses the same status/file-line dialect as batch; no multi-count
  summary is required on single (batch-only).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.UseMany = false
	return nil
}
```
