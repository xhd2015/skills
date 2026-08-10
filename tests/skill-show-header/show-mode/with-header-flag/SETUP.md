# Scenario

**Feature**: `skill --show --header` prints YAML frontmatter only

```
# header-only invocation
user -> go-best-practice skill --show --header -> ---\nname: ...\n---
```

## Preconditions

- `req.HeaderOnly` is true; flag order is `--show` then `--header`.

## Steps

1. Set `req.HeaderOnly = true`.

## Context

- stdout must include delimiter lines and `name:` but must omit `# Go Best Practice Skill`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.HeaderOnly = true
	req.HeaderBeforeShow = false
	return nil
}
```
