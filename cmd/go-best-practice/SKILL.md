---
name: go-best-practice
description: >-
  Index of Go best-practice recipes (project scaffolding, CLI flag
  parsing, and more). Use when the user wants to bootstrap a new
  project or parse CLI flags in a Go program. Load a sub-topic with:
  go-best-practice skill --show <topic-path>
---

# Go Best Practice Skill

This skill is an **index**. Load a detailed recipe with
`go-best-practice skill --show <topic>` (or
`go-best-practice skill <topic> --show`). Topics are organized as a
tree; address a sub-topic with a slash-separated path, e.g.
`flags-parsing/types` or `cli/color`.

## Topics

- `kool-create` — scaffold new projects with `kool create` (react,
  go-react, frontend, server, electron)
- `cmd-exec` — running external commands with
  `github.com/xhd2015/xgo/support/cmd` (Debug mode, output capture,
  env vars, directory, I/O redirect)
- `cli` — CLI UX and skill CLI packaging
  - `color` — terminal ANSI color: `--color` / `--no-color`, TTY
    auto, and the `NO_COLOR` env convention
  - `streaming` — stream CLI output as work proceeds; avoid
    buffering all output until the end (when to buffer, flush,
    NDJSON vs full JSON)
  - `skill-cli` — skill CLI shapes: single-skill, multi-skill host,
    topic discovery
- `flags-parsing` — CLI flag parsing with
  `github.com/xhd2015/less-flags`
  - `types` — supported target types (`*bool`, `*string`, `*int`,
    `*time.Duration`, `*[]string`, `Cut`, and `**T` variants)
  - `subcommand` — sub-command dispatcher patterns (with `StopOnFirstArg` and no-toplevel-flags variants)
  - `cut` — cut flags: consume all remaining tokens after a marker
  - `collect` — `CollectParsedFlags` / `Flags.Reconstruct` / `Remove`

## Usage

```bash
# list top-level topics (CLI help index)
go-best-practice
go-best-practice topics

# root skill index
go-best-practice skill --show

# reveal a top-level topic
go-best-practice skill --show kool-create
go-best-practice skill --show flags-parsing
go-best-practice skill --show cli

# reveal a sub-topic (slash-separated path; both flag orders)
go-best-practice skill --show cli/color
go-best-practice skill --show cli/streaming
go-best-practice skill --show cli/skill-cli
go-best-practice skill --show flags-parsing/types
go-best-practice skill flags-parsing/subcommand --show
go-best-practice skill --show flags-parsing/cut
go-best-practice skill --show flags-parsing/collect
```
