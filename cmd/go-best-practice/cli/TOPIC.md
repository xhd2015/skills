---
name: go-best-practice/cli
description: >-
  CLI UX (color, streaming output) and skill CLI packaging shapes.
  Load a child with: go-best-practice skill --show cli/<topic>
---

# cli — CLI UX and skill CLI packaging

Recipes for building Go CLIs: how output looks and streams, and how to
ship skill binaries that embed `SKILL.md` / nested `TOPIC.md` trees.

This is a **category index**. `color` and `streaming` are general CLI
I/O UX; `skill-cli` is how to package skill CLIs. Flag parsing lives
separately under `flags-parsing`.

## Topics

- `color` — terminal ANSI color: `--color` / `--no-color`, TTY auto,
  and the `NO_COLOR` env convention
- `streaming` — stream CLI output as work proceeds; avoid buffering
  all output until the end (when to buffer, flush, NDJSON vs full JSON)
- `skill-cli` — skill CLI shapes: single-skill, multi-skill host,
  topic discovery

## Retrieve

```bash
go-best-practice skill --show cli
go-best-practice skill --show cli/color
go-best-practice skill --show cli/streaming
go-best-practice skill --show cli/skill-cli
go-best-practice skill cli/color --show
```

## See also

- `flags-parsing` — less-flags and sub-command `--help` at every level
- `cmd-exec` — running external commands (inherit stdout/stderr)
