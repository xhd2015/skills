# Scenario

**Feature**: folded description block parses to one normalized string

```
# multi-line folded description in header
SKILL.md content -> GetHeader -> folded YAML

# parser normalizes folded value
inner YAML -> ParseHeader -> Entries.Get("description")
```

## Preconditions

- Fixture uses `description: >-` followed by indented continuation lines.

## Steps

1. Load `skill_content.md` into `req.Content`.

## Context

- Expected normalized description is `multi line`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Content = readFixture(t, d, "skill_content.md")
	return nil
}
```