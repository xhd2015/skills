# Scenario

**Feature**: install writes nested skill-cli/TOPIC.md extras (not nested SKILL.md / topics/*)

```
# TreeFS ExtraFiles derived from **/TOPIC.md
caller -> SingleSkill.Handle(--install)
  -> .agents/skills/demo-skill/SKILL.md
  -> .agents/skills/demo-skill/skill-cli/TOPIC.md
  -> no skill-cli/SKILL.md, no topics/skill-cli.md
```

## Preconditions

- Tree SingleSkill configured by parent with TOPIC.md tree files.
- ExtraFiles left nil so install derives nested TOPIC.md paths from TreeFS.

## Steps

1. Enable workdir and install without dry-run.
2. Do not override ExtraFiles (harness derives from TreeFiles TOPIC.md paths).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.UseWorkDir = true
	req.Args = []string{"--install"}
	// ExtraFiles remains nil: buildSingleSkill / SingleSkill collect TOPIC.md only.
	return nil
}
```
