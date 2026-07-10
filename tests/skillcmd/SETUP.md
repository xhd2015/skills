# Scenario

**Feature**: skillcmd is the shared framework for skill CLI flag actions

```
# parse classifies skill argv into one action + rest
caller -> skillcmd.ParseSkillArgs(args) -> Action + Header + Rest

# single / registry handlers drive show, list, install, update
caller -> SingleSkill.Handle / Registry.HandleSkill -> stdout | filesystem
caller -> skillcmd.HandleInstall / HandleUpdateMany -> .agents/skills/<name>/

# nested multi-topic TreeFS uses path/TOPIC.md (root stays SKILL.md)
caller --show a/b -> loadTreeSkill -> a/b/TOPIC.md
caller --list -> ListTreeTopics(**/TOPIC.md)

# header helpers wrap YAML frontmatter
caller -> GetHeader / FormatHeaderWithDelimiters -> YAML text
```

## Preconditions

- Package `github.com/xhd2015/skills/skillcmd` implements parse, single, registry,
  file header, and install/update APIs used by `Run`.
- Leaves set `req.Mode` and mode-specific fields; install/update leaves use
  isolated workdirs and optional `HOME`.

## Steps

1. Leaves assign `req.Mode` and scenario inputs (Args, skill content, registry).
2. `Run` dispatches to the matching skillcmd API and captures stdout/errors.
3. Assert checks structured parse results or stdout/filesystem effects.

## Context

- Canonical demo skill content uses name `demo-skill` and body marker
  `# Demo Skill Body` unless a leaf overrides.
- Nested tree demos use `a/b/TOPIC.md` and `skill-cli/TOPIC.md` paths
  (not nested `SKILL.md` / `topics/*.md`).

```go
import "testing"

const (
	demoSkillName = "demo-skill"
	demoRootContent = `---
name: demo-skill
description: demo skill for skillcmd doctests
---

# Demo Skill Body

Root index content.
`
	nestedABContent = `---
name: demo-skill/a/b
description: nested topic a/b
---

# Nested A/B Body

Nested topic content for a/b.
`
	nestedSkillCLIContent = `---
name: demo-skill/skill-cli
description: nested skill-cli topic
---

# Nested Skill-CLI Body

skill-cli nested TOPIC.md.
`
	fooSkillContent = `---
name: foo
description: foo skill description
---

# Foo Skill Body
`
	barSkillContent = `---
name: bar
description: bar skill description
---

# Bar Skill Body
`
)

func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	if req.TreeFiles == nil {
		req.TreeFiles = map[string]string{}
	}
	if req.PostInstallMutate == nil {
		req.PostInstallMutate = map[string]string{}
	}
	return nil
}
```
