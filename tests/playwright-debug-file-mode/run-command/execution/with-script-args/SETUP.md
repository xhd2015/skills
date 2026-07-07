# Scenario

**Feature**: `run <file.js>` forwards trailing args to `process.argv`

```
user -> playwright-debug CLI (run print-argv.js <scriptArgs...>) -> JSON argv on stdout
```

## Preconditions

- Requires `node`, `npm`, and playwright cache (first run may download Chromium).

## Steps

1. Each leaf sets `req.Args` with `run`, an absolute `print-argv.js` path, and trailing script args.

## Context

- All leaves in this group are labeled `slow`.
- stdout is `JSON.stringify(process.argv.slice(3))` from the user script.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```