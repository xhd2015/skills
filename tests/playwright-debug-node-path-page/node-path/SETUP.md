# Scenario

**Feature**: NODE_PATH points node subprocess at playwright cache node_modules

```
# CLI sets NODE_PATH before spawning node for file mode
user -> playwright-debug CLI (run <fixture.js>) -> node + NODE_PATH -> nested require('playwright') resolves
```

## Preconditions

- NODE_PATH must be `<cacheDir>/node_modules` where cacheDir is
  `~/.playwright-debug/node_package`.
- Nested modules cannot rely on a local node_modules beside the fixture.

## Steps

1. Each leaf runs a fixture that exercises NODE_PATH resolution or visibility.

## Context

- All leaves in this group are labeled `slow`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```