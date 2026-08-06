# Scenario

**Feature**: SingleSkill with TreeFS resolves nested path/TOPIC.md

```
# nested topic show, list, and install extra files use TOPIC.md
caller -> SingleSkill.Handle(--show a/b) -> a/b/TOPIC.md content
caller -> SingleSkill.Handle(--list) -> topics from **/TOPIC.md
caller -> SingleSkill.Handle(--install) -> ExtraFiles include skill-cli/TOPIC.md
```

## Preconditions

- TreeFiles map provides nested TOPIC.md paths under the skill root.
- Root index content stays on SingleSkill.RootContent (SKILL.md at install root).

## Steps

1. Set `req.Mode = ModeSingle` with demo root + tree files.
2. Leaves set show/list/install Args and optional ExtraFiles overrides.

## Context

- Invalid path segments (`.`, `..`, empty) are rejected.
- Nested `SKILL.md` is not the discoverable topic filename.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Mode = ModeSingle
	req.SkillName = demoSkillName
	req.RootContent = demoRootContent
	req.Usage = "demo-skill skill --install"
	req.TreeFiles = map[string]string{
		"a/b/TOPIC.md":       nestedABContent,
		"skill-cli/TOPIC.md": nestedSkillCLIContent,
	}
	return nil
}
```
