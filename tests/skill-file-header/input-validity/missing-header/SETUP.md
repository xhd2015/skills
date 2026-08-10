# Scenario

**Feature**: content without opening frontmatter delimiter is rejected

```
# body-only SKILL.md has no YAML block
SKILL.md content -> GetHeader -> error
```

## Preconditions

- Fixture `skill_content.md` does not start with `---\n`.

## Steps

1. Load `skill_content.md` into `req.Content`.

## Context

- `ParseHeader` must not run when `GetHeader` fails.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Content = readFixture(t, d, "skill_content.md")
	return nil
}
```