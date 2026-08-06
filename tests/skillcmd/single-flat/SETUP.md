# Scenario

**Feature**: SingleSkill without TreeFS supports show, list, and install

```
# flat single skill host
caller -> SingleSkill.Handle(--show|--list|--install...) -> stdout / install layout
```

## Preconditions

- Skill has Name + RootContent only (no nested tree).

## Steps

1. Set `req.Mode = ModeSingle`.
2. Configure demo skill identity and root content.
3. Leaves set action Args.

## Context

- Body marker `# Demo Skill Body` distinguishes full show from header-only.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Mode = ModeSingle
	req.SkillName = demoSkillName
	req.RootContent = demoRootContent
	req.Usage = "demo-skill skill --install"
	return nil
}
```
