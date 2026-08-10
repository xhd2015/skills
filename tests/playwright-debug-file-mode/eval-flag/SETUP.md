# Scenario

**Feature**: `-e` / `--eval` explicit adhoc eval mode

```
user -> playwright-debug CLI (-e|--eval <script>) -> eval runner -> stdout marker
```

## Preconditions

- Eval flag requires a script argument; extra positional args are rejected.

## Steps

1. Error leaves test malformed flag usage.
2. Execution leaf runs minimal `console.log` eval snippet.

## Context

- Routing errors are fast; execution is labeled `slow`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```