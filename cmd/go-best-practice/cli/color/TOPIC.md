---
name: go-best-practice/cli/color
description: >-
  Terminal ANSI color for CLI output: --color / --no-color, TTY auto,
  and the NO_COLOR environment convention.
---

# color — terminal ANSI color for CLIs

Use `github.com/xhd2015/dot-pkgs/go-pkgs/terminal/color`. Do not copy
`Resolve`, `Style`, or TTY helpers.

When a Go CLI prints human-facing status, summaries, or highlights,
gate ANSI with the library's three-mode policy. Default is **auto**.
Force-on is only `--color`. Do **not** implement `FORCE_COLOR`,
`CLICOLOR`, or `CLICOLOR_FORCE`.

## Policy

| Mode | How selected | Color on? |
| ---- | ------------ | --------- |
| **Auto** (default) | neither flag | TTY → on, unless `NO_COLOR` non-empty → off; non-TTY → off |
| **Always** | `--color` | always on (ignores TTY and `NO_COLOR`) |
| **Never** | `--no-color` | always off (ignores TTY and `NO_COLOR`) |

**Conflict:** `--color` and `--no-color` together must fail with:

```text
--color and --no-color cannot be specified together
```

Flags always win. `NO_COLOR` applies only in Auto: unset or empty has no
effect; any non-empty value (`NO_COLOR=1`) disables. See
[no-color.org](https://no-color.org/).

## Use the library

Parse `--color` / `--no-color` as bools with your flag library
(`less-flags` `Bool`). Then `ModeFromFlags`. Do not hand-scan argv.
Color against the writer you print to (stdout vs stderr).

```go
import "github.com/xhd2015/dot-pkgs/go-pkgs/terminal/color"

mode, err := color.ModeFromFlags(colorFlag, noColorFlag)
if err != nil {
    return err
}
style := color.Style{Enabled: color.EnabledFor(mode, stdout)}
fmt.Fprintln(stdout, style.Green("ok"))
fmt.Fprintf(stderr, "%s skipped 3\n", style.Yellow("warning:"))
```

Help text:

```text
  --color      force ANSI color on (even when stdout is not a TTY)
  --no-color   force ANSI color off
```

## Tokens

- **Yellow** `warning:` on stderr
- **Red** `Error:` on stderr
- **Green** success / status on stdout
- **Gray** meta (duration, counts) on stdout
- **No ANSI** in JSON or other machine-readable output

## Tests

Harness stdout is usually a **pipe**, so Auto is off unless `--color`.
Library tests inject `noColorEnv` into `Resolve` — do not `t.Setenv`.
CLI tests force `--color` / `--no-color`. Assert specific SGR sequences
when useful (e.g. green pass), not only “any escape present”.

| Scenario | Setup | Expect |
| -------- | ----- | ------ |
| Force on | `--color` | ANSI present |
| Force off | `--no-color` | no ANSI |
| Force wins | `--color` + `NO_COLOR=1` | ANSI present |
| Auto pipe | no flags, stdout pipe | no ANSI |
| Conflict | both flags | error: `cannot be specified together` |

## Out of scope

- `FORCE_COLOR`, `CLICOLOR`, `CLICOLOR_FORCE`
- Full terminfo / 256-color palettes
- Coloring JSON or machine-readable output

## See also

- `cli/streaming` — stream output as work proceeds
- `cli/config` — gray `notice:` when a preference comes from config.json
- `flags-parsing` — bool flags and help text with less-flags
- [no-color.org](https://no-color.org/) — `NO_COLOR` convention

Reveal with:

```bash
go-best-practice skill --show cli/color
```
