# Scenario

**Feature**: `-e` / `--eval` forwards trailing args to `process.argv`

```
user -> playwright-debug CLI (-e|--eval <script> <scriptArgs...>) -> JSON argv on stdout
```

## Preconditions

- Requires playwright cache (slow).

## Steps

1. Each leaf sets eval flag, inline script, and trailing script args on `req.Args`.

## Context

- All leaves labeled `slow`.
- Inline script logs `JSON.stringify(process.argv.slice(3))`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```