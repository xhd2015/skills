---
name: go-best-practice/cli
description: >-
  CLI UX (color, streaming output, dry-run, inline TUI mouse) and skill
  CLI packaging shapes. Load a child with: go-best-practice skill
  --show cli/<topic>
---

# cli — CLI UX and skill CLI packaging

Recipes for building Go CLIs: how output looks and streams, interactive
terminal UIs, persisted preferences, and how to ship skill binaries that
embed `SKILL.md` / nested `TOPIC.md` trees.

This is a **category index**. `color`, `streaming`, `dry-run`, and
`config` are general CLI I/O UX; `inline-tui-mouse` is mouse hit-testing
for inline TUIs; `skill-cli` is how to package skill CLIs. Flag parsing
lives separately under `flags-parsing`.

## Topics

- `color` — terminal ANSI color via `terminal/color`: `--color` /
  `--no-color`, TTY auto, and the `NO_COLOR` env convention
- `streaming` — stream CLI output as work proceeds; avoid buffering
  all output until the end (when to buffer, flush, NDJSON vs full JSON)
- `dry-run` — one pipeline with side-effect gates; avoid a separate
  dry-run function that duplicates logic
- `config` — persist flag preferences in tool-home `config.json`:
  `--set-config` / `--show-config` / `--no-config`, precedence, gray
  `notice:` when a value comes from config
- `skill-cli` — skill CLI shapes: single-skill, multi-skill host,
  topic discovery
- `inline-tui-mouse` — mouse hit-testing for inline (non-alt-screen)
  TUIs: view-local hitmaps, CSI 6n origin on one stdin path, dual-origin
  fallback, anti-patterns

## Retrieve

```bash
go-best-practice skill --show cli
go-best-practice skill --show cli/color
go-best-practice skill --show cli/streaming
go-best-practice skill --show cli/dry-run
go-best-practice skill --show cli/config
go-best-practice skill --show cli/skill-cli
go-best-practice skill --show cli/inline-tui-mouse
go-best-practice skill cli/color --show
```

## See also

- `flags-parsing` — less-flags and sub-command `--help` at every level
- `cmd-exec` — running external commands (inherit stdout/stderr)
