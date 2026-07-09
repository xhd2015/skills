# Scenario

**Feature**: `skill --header --show` is accepted (both flag orders)

```
# header flag may appear before --show
user -> go-best-practice skill --header --show -> YAML frontmatter only
```

## Preconditions

- Header-only mode with HeaderBeforeShow.

## Steps

1. Set HeaderOnly and HeaderBeforeShow true.

## Context

- Output must match `--show --header` (delimiters + name, no body marker).

```go
func Setup(t *testing.T, req *Request) error {
	req.HeaderOnly = true
	req.HeaderBeforeShow = true
	return nil
}
```
