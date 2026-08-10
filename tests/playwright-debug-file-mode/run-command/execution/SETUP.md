# Scenario

**Feature**: `run <file.js>` executes scripts via bootstrap file mode

```
user -> playwright-debug CLI (run <fixture.js>) -> file runner -> marker on stdout
```

## Preconditions

- Requires `node`, `npm`, and playwright cache (first run may download Chromium).

## Steps

1. Each leaf sets `req.Args = []string{"run", <absolute fixture path>}`.

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