# Scenario

**Feature**: valid SKILL.md returns inner YAML and parsed name/description entries

```
# standard two-field header
SKILL.md content -> GetHeader -> name + description YAML

# entries support lookup
inner YAML -> ParseHeader -> Entries.Get("name")
```

## Preconditions

- Fixture `skill_content.md` contains `name: git-fetch` and `description: clone repos`.

## Steps

1. Load `skill_content.md` into `req.Content`.

## Context

- Delimiters must not appear in the `GetHeader` return value.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Content = readFixture(t, d, "skill_content.md")
	return nil
}
```