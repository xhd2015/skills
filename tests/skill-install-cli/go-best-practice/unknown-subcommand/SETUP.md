# Scenario

**Bug**: unknown `skill` sub-commands should list valid choices including `skill install`

```
# bogus skill sub-command rejected with helpful message
user -> go-best-practice skill bogus -> unknown skill sub-command error
```

## Preconditions

- `go-best-practice` binary is built and on `req.Binary`.

## Steps

1. Set `req.Args = ["skill", "bogus"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "bogus"}
	return nil
}
```