# Go Best Practice Review — `github.com/xhd2015/skills`

**Date:** 2026-08-06  
**Scope:** codebase structure, CLI design, flag handling, package layout  
**Lens:** in-repo `go-best-practice` topics (`cli/*`, `flags-parsing/*`, `cmd-exec`, `go-embed-assets`, `kool-create`, `cli/skill-cli`)  
**Mode:** review only — no behavioral fixes applied

---

## Project snapshot

| Area | What exists |
|------|-------------|
| Module | `github.com/xhd2015/skills` (Go 1.25) |
| Binaries | `cmd/go-best-practice`, `cmd/github-fetch`, `cmd/playwright-debug` |
| Libraries | `skillcmd` (skill CLI home), `playwrightdebug`, deprecated shims `install` / `skill_file` |
| Tests | Unit tests under packages + large `tests/` doctest trees |
| Flag deps | **Two** parsers: `github.com/xhd2015/less-flags` and `github.com/xhd2015/less-gen/flags` |

**Shape summary (skill-cli):**

| CLI | Shape | Notes |
|-----|-------|--------|
| `go-best-practice` | Shape 3 (topic tree) | `SingleSkill` + `TreeFS`; good model host |
| `github-fetch` | Shape 1 + domain cmds | Many word subcommands; help uneven |
| `playwright-debug` | Shape 1 + domain cmds | Launch flags peeled manually; eval path heavy |

---

## Findings (by severity)

### High

#### H1. Dual flag libraries split the codebase

**Topic:** `flags-parsing` (less-flags is the taught stack)

| Location | Import |
|----------|--------|
| `cmd/github-fetch/*`, `cmd/go-best-practice/vet` | `github.com/xhd2015/less-flags` |
| `skillcmd/install.go`, `skillcmd/update.go` | `github.com/xhd2015/less-gen/flags` |

The skill **recipes** and most domain CLIs standardize on `less-flags`. Shared install/update code — the library every skill CLI goes through — uses a different package (`less-gen/flags`). That means:

- Two mental models / APIs for the same job  
- Divergent error types (`lessflags.ErrHelp` vs `flags.ErrHelp`) already appear in code  
- Future flag features (Cut, CollectParsedFlags, types docs) won’t apply evenly  

**Recommended change:** migrate `skillcmd` install/update parsing to `github.com/xhd2015/less-flags` only; drop `less-gen` if nothing else needs it. Keep one import path module-wide.

---

#### H2. Missing `--help` at several command levels (`github-fetch`)

**Topics:** `flags-parsing/subcommand`, `cli` (“every level needs `--help`”)

| Command | Current behavior | Gap |
|---------|------------------|-----|
| `work-on` | No flags, no help; bare `work-on` → error only | No `work-on --help` / empty-args help |
| `prs` / `pulls` / `issues` | `parseListFlags` has **no** `Help(...)` | `prs --help` fails as unknown flag or never shows usage |
| `issue` | Only `--list`; no help text | Dead end for explorers |
| Root `help` | Omits `prs`, `issues`, aliases | Discoverability broken vs switch table |

`pr`, `ci`, `push`, `status`, `yaml` largely do the right thing (less-flags `Help` or manual help). Root help also omits the recipe line:

> `Run github-fetch <command> --help for command-specific options.`

**Recommended changes:**

1. Add `listHelp` (and work-on help) with `lessflags.Help("-h,--help", …)` in `parseListFlags` / `handleWorkon`.  
2. Expand root `help` to list **all** dispatch names (`prs`/`pulls`, `issues`, `issue`, `fetch`, `checks`, …).  
3. Point users to per-command `--help` from the root text.

---

#### H3. README / docs disagree with real CLI surface

**Topics:** `cli/skill-cli` (action flags, not word subcommands), Shape 3 retrieve path

| Doc claim | Actual CLI |
|-----------|------------|
| `go-best-practice flags-parsing` | **Rejected** — must be `skill --show flags-parsing` (or `topics`) |
| `playwright-debug skill install --cursor` | Invalid — install is **`--install`**, not a word subcommand |
| README omits `github-fetch` entirely | Third major binary with no install examples |

Agents and humans following README will hit errors. `cli/skill-cli` is explicit: actions are flags only (`--show` / `--install` / `--list`).

**Recommended changes:** fix README examples to real invocations; document `github-fetch`; optionally add a **compatibility** bare-topic path for `go-best-practice` only if product wants that UX (today code intentionally errors).

---

#### H4. `github-fetch` package layout: oversized `package main` monorepo

**Topics:** package layout / maintainability; `vet`’s own `--file-max-lines` default (500)

| File | ~LOC | Role |
|------|------|------|
| `cmd/github-fetch/github_api.go` | 780 | HTTP, git, list flags, PR types, helpers |
| `cmd/github-fetch/main_test.go` | 1890 | All CLI tests in one file |
| `cmd/github-fetch/ci.go` | 454 | CI UX + log assembly |

Everything is `package main`. Domain logic (API client, auth, list formatting, git) is not reusable and hard to unit-test without CLI entrypoints. This is the opposite of `playwright-debug` (thin `cmd/` + `playwrightdebug` library).

**Recommended changes:**

```text
cmd/github-fetch/          # main, dispatch, help only
githubfetch/               # or internal/githubfetch/
  api/                     # HTTP + types
  auth/
  list/
  ci/
  gitutil/                 # runGitCmd wrappers
```

Split tests by command area. Align with `vet --file-max-lines 500` for production `.go` files.

---

### Medium

#### M1. Manual launch-flag peeling in `playwright-debug`

**Topics:** `flags-parsing`, `flags-parsing/types` (StringSlice, Bool)

`playwrightdebug.ExtractLaunchFlags` hand-parses `--extension`, `--load-extension`, `--user-data-dir`, `--headed`/`--headless`, plus `=` forms. That reimplements less-flags and is easy to drift (e.g. incomplete `=` support for some flags, no unified help surface for launch flags alone).

**Recommended change:** peel with less-flags where possible, e.g.:

- `StringSlice` for `--extension` / `--load-extension`  
- `String` for `--user-data-dir`  
- mutually exclusive headed/headless bools  
- `Help` block for launch options  

If launch flags must interleave with script argv / cut semantics, document the peel order and still centralize parsing (or use `Cut` for trailing script args after a marker).

---

#### M2. External commands use raw `os/exec` everywhere

**Topic:** `cmd-exec` (`github.com/xhd2015/xgo/support/cmd`)

Call sites include:

| Area | Commands |
|------|----------|
| `github-fetch` | `git`, `gh`, `actionlint` via `runCmd` / `exec.Command` |
| `playwrightdebug` | `node`, `npm` |
| `cmd/playwright-debug` | `node -e` eval wrapper |
| `go-best-practice/vet` | `git rev-parse` |

Recipe preference: fluent `cmd.Debug().Dir(...).Env(...).Run(...)` / `cmd.Output(...)` for visibility and consistent I/O inheritance.

**Recommended change:** introduce thin wrappers that use `xgo/support/cmd` for user-visible invocations (`git push`, `actionlint`, `npm install`, …). Keep silent capture paths on `cmd.Output` or non-Debug `Run`. Not mandatory for every test helper, but production CLIs should converge.

---

#### M3. CI log path buffers entire multi-job output (`cli/streaming`)

**Topic:** `cli/streaming`

In `showRunLogs` (`ci.go`), a `strings.Builder` accumulates workflow header + every job’s truncated log, then one `WriteString` at the end. For multi-job runs this:

- Delays first bytes until all jobs are fetched  
- Holds full log text in memory  

Streaming is the default design; buffering is justified for tables/sort/atomic JSON — not for sequential log dumps.

**Recommended change:** write header immediately; for each job, write job header + truncated log as soon as `fetchJobLogs` returns. Keep truncation logic; drop the full-buffer pattern.

---

#### M4. Eval mode duplicates bootstrap launch logic

**Topics:** package layout, `go-embed-assets` (single embed source of truth), DRY

File mode correctly embeds `bootstrap.cjs` and runs it. Eval mode in `cmd/playwright-debug/main.go` inlines a large JS `launchBrowser()` (~100 lines) that must stay in sync with bootstrap (comment already says “keep in sync”).

**Recommended change:** eval should also use embedded bootstrap (e.g. temp script that `require`s bootstrap helpers, or write snippet to temp file and `RunFile`). One launch implementation.

---

#### M5. Install `--force` is parsed but not documented in `--help`

**Topics:** `flags-parsing` (help next to flags), `cli/skill-cli` install flags

`HandleInstall` binds `--force` and clears `noOverride`, but the generated help string lists `--no-override` and omits `--force`. Users discover force only via tests/behavior.

**Recommended change:** document `--force` in install help (and note interaction with `--no-override`).

---

#### M6. `go-best-practice` main still owns parallel topic helpers

**Topics:** `cli/skill-cli` Shape 3, package simplicity

`SingleSkill{TreeFS: skillTreeFS}` already implements show/list/install for nested `TOPIC.md`. Yet `main.go` still defines `readTopic`, `validateSegments`, `collectTopicFiles` (used primarily by unit tests, not by the live install path which goes through `skillcmd`).

**Recommended change:** delete or move helpers into tests only; call `skillcmd.ListTreeTopics` / public show path for assertions. Shrink `main.go` to dispatch + embed wiring.

---

#### M7. `cli/skill-cli` recipe still shows deprecated packages as primary

**Topics:** `cli/skill-cli` (self-consistency of the skill that documents the monorepo)

`skillcmd` is correctly documented as the home, but **brief main.go** samples still import:

```go
"github.com/xhd2015/skills/install"
"github.com/xhd2015/skills/skill_file"
```

instead of `skillcmd.SingleSkill` / `HandleInstall` / header helpers. New CLIs copied from the recipe will reintroduce deprecated packages (and dual APIs). Doctests still exercise `install` shims — fine for compatibility — but the **canonical** sample should be `skillcmd` only.

**Recommended change:** rewrite Shape 1/2/3 samples around `skillcmd.SingleSkill` / `Registry` (as `cmd/go-best-practice` already does). Keep a short “deprecated shims” note.

---

#### M8. Color policy only on update path

**Topic:** `cli/color`

`skillcmd` implements `ColorMode` / `ResolveColor` / `NO_COLOR` correctly for **update** summaries. Install and domain CLIs (`github-fetch status`, CI headers, install plans) either always plain or ad-hoc.

**Recommended change:** when adding ANSI to any human status path, reuse `skillcmd` color helpers + `--color` / `--no-color` conflict rule. Do not invent a second color policy.

---

### Low

#### L1. Incomplete root / skill docs for `github-fetch` domain surface

`SKILL.md` covers `pr`, `work-on`, `push` only. Live CLI also has `ci`, `status`, `yaml validate`, list commands, skill actions. Agents installing the skill miss half the product.

**Recommended change:** extend SKILL.md with domain workflows (not install plumbing). Keep install in CLI help / README only (`cli/skill-cli` SKILL.md rules).

---

#### L2. `status` help only if `-h` is the first remaining arg

**Topic:** `flags-parsing` (prefer less-flags Help)

`handleStatus` checks `args[0] == -h|--help`. Prefer `lessflags.Help` (or scan consistently) so `status --help` is the only documented shape and unknown flags error cleanly.

---

#### L3. `yaml validate` hand-scans argv for help

Works, but inconsistent with less-flags elsewhere. Use less-flags `Help` + positional file for one pattern.

---

#### L4. `issue` is a thin alias with poor UX

`issue --list` redirects to `issues`; without `--list` it errors. Prefer documenting only `issues`, or make bare `issue` show help / list open issues.

---

#### L5. `kool-create` not used here (N/A, not a defect)

This monorepo is not a greenfield app scaffold. No change required. Topic remains for consumers of the skill, not for restructuring this module.

---

#### L6. `go-embed-assets` largely N/A for fat SPA trees; current embeds are fine

`//go:embed SKILL.md`, topic dirs, and `bootstrap.cjs` are compile-safe and tracked. No empty-embed anti-pattern. Full Layer 4 hydrate is unnecessary for markdown/bootstrap payloads.

---

## What already matches best practice

These are strengths to preserve:

1. **`skillcmd` as single skill-CLI home** — install inventory is plan-then-apply (`planInventory` → gate dry-run → `apply`), matching `cli/dry-run`.  
2. **Shape 3 on `go-best-practice`** — `TreeFS`, topic index on help/list, both `--show` orders via `ParseSkillArgs`.  
3. **Color implementation** on update — three modes, conflict error text, `NO_COLOR` only in auto (`cli/color`).  
4. **Domain less-flags usage** on `pr` / `ci` / `push` / `vet` with `Help` + often `HelpNoExit`.  
5. **Playwright library split** — `cmd/playwright-debug` thin; `playwrightdebug` owns ensure/run/launch.  
6. **Deprecated shims** — `install` / `skill_file` re-export `skillcmd` instead of forking logic.  
7. **Streaming-friendly lists** — PR/issue list printers write line-by-line, not full-buffer join.  
8. **Skill SKILL.md hygiene tests** — `github-fetch` tests reject install plumbing phrases in skill body.

---

## Recommended change backlog (priority order)

| Priority | Item | Primary topics |
|----------|------|----------------|
| 1 | Docs: fix README invocations; document `github-fetch` | skill-cli, CLI UX |
| 2 | Docs + code: full command inventory + `--help` on `work-on`, list cmds | flags-parsing/subcommand |
| 3 | Docs: install `--force` in help | flags-parsing |
| 4 | Unify on `less-flags` in `skillcmd` | flags-parsing |
| 5 | Stream CI logs job-by-job | cli/streaming |
| 6 | Deduplicate playwright eval vs bootstrap | package layout, embed |
| 7 | less-flags for launch peel (or Cut design) | flags-parsing |
| 8 | Split `github-fetch` out of fat `package main` | package layout |
| 9 | Prefer `xgo/support/cmd` for external tools | cmd-exec |
| 10 | Rewrite skill-cli samples to `skillcmd` only | cli/skill-cli |
| 11 | Drop dead topic helpers from `go-best-practice` main | skill-cli Shape 3 |
| 12 | Expand `github-fetch` SKILL.md domain coverage | skill-cli SKILL rules |

Items 1–3 are high leverage and low risk (docs + help wiring). Items 4–8 are structural. Items 9–12 are consistency polish.

---

## Package layout target (suggested)

```text
github.com/xhd2015/skills/
├── cmd/
│   ├── go-best-practice/     # Shape 3 host + topics + vet
│   ├── github-fetch/         # thin main only
│   └── playwright-debug/     # thin main only
├── skillcmd/                 # parse, single, registry, install, update, color
├── playwrightdebug/          # already good
├── githubfetch/              # NEW: extract from cmd/github-fetch
├── install/, skill_file/     # keep as deprecated shims until callers migrate
└── tests/                    # doctests (migrate imports to skillcmd over time)
```

---

## Applicability matrix

| Topic | Applicable? | Verdict |
|-------|-------------|---------|
| `flags-parsing` (+ types/subcommand/cut/collect) | Yes | Partial — dual lib; help gaps; manual peel; cut/collect unused (OK if no need) |
| `cli/skill-cli` | Yes | Strong runtime; recipe samples lag `skillcmd` |
| `cli/dry-run` | Yes (install/update) | Good plan-then-apply |
| `cli/color` | Yes | Good on update only |
| `cli/streaming` | Yes | Lists good; CI logs buffer |
| `cli/inline-tui-mouse` | No | No inline TUI |
| `cmd-exec` | Yes | Not adopted; raw `os/exec` |
| `go-embed-assets` | Partial | Correct for current embeds; hydrate N/A |
| `kool-create` | No (scaffold other projects) | Out of scope |

---

## Conclusion

The monorepo has a solid **skillcmd** core and a well-shaped **go-best-practice** topic tree. The largest gaps versus its own recipes are: **inconsistent flag stacks**, **incomplete per-command help and root/README surfaces**, **fat `github-fetch` main package**, and **playwright eval / streaming / external-command patterns** that lag the documented standards.

No code fixes were made in this pass. Prefer implementing the backlog from docs/help (H2–H3, M5) before structural refactors (H1, H4, M4, M7).
