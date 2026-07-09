# Scenario

**Feature**: SingleSkill with TreeFS resolves nested path/SKILL.md

```
# nested topic show and install extra files
caller -> SingleSkill.Handle(--show a/b) -> a/b/SKILL.md content
caller -> SingleSkill.Handle(--install) -> ExtraFiles include skill-cli/SKILL.md
```

## Preconditions

- TreeFiles map provides nested SKILL.md paths under the skill root.

## Steps

1. Set `req.Mode = ModeSingle` with demo root + tree files.
2. Leaves set show/install Args and optional ExtraFiles.

## Context

- Invalid path segments (`.`, `..`, empty) are rejected.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = ModeSingle
	req.SkillName = demoSkillName
	req.RootContent = demoRootContent
	req.Usage = "demo-skill skill --install"
	req.TreeFiles = map[string]string{
		"a/b/SKILL.md":        nestedABContent,
		"skill-cli/SKILL.md": nestedSkillCLIContent,
	}
	return nil
}
```
