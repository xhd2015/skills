# Scenario

**Feature**: `skill --show` without flags prints full SKILL.md content

```
# default show invocation
user -> go-best-practice skill --show -> header + Markdown body
```

## Preconditions

- `req.HeaderOnly` is false.

## Steps

1. Set `req.HeaderOnly = false`.

## Context

- stdout must include both YAML header fields and the body marker.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.HeaderOnly = false
	req.HeaderBeforeShow = false
	return nil
}
```
