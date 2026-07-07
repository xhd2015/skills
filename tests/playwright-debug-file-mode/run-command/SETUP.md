# Scenario

**Feature**: explicit `run <file.js>` file mode

```
# router requires existing .js file for run subcommand
user -> playwright-debug CLI (run <path>) -> file runner | routing error
```

## Preconditions

- `run` never falls back to eval for non-file strings.

## Steps

1. Error leaves pass invalid `run` arg combinations.
2. Execution leaves pass absolute paths to shared `testdata/` fixtures.

## Context

- Error leaves are fast (no browser). Execution leaves are labeled `slow`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```