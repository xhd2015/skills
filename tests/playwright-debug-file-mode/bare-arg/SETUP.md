# Scenario

**Feature**: single bare positional arg routes to file or eval mode

```
# existing .js path → file alias
user -> playwright-debug CLI (<file.js>) -> file runner

# non-file string → eval shorthand
user -> playwright-debug CLI ('<script>') -> eval runner
```

## Preconditions

- File alias only applies when the single arg is an existing `.js` file.

## Steps

1. File-alias leaf passes absolute path to `simple-eval.js` without `run`.
2. Eval leaf passes inline script string without `-e`.

## Context

- Both leaves are labeled `slow`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```