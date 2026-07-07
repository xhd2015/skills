# Scenario

**Feature**: nested library functions receive `page` as an explicit parameter

```
# bootstrap injects page; entry script passes it into nested lib
user -> playwright-debug CLI (run <fixture.js>) -> greet(page) -> page.goto -> explicit-page-ok
```

## Preconditions

- Nested CommonJS modules run in isolated scope without implicit global `page`.

## Steps

1. Run fixture that requires a nested lib and passes injected `page`.

## Context

- Leaf is labeled `slow` (launches browser).

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```