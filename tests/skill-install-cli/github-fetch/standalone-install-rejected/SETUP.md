# Scenario

**Feature**: github-fetch rejects top-level `install` (must use `skill --install`)

```
# standalone install is not a registered command
user -> github-fetch install -> unknown command error
```

## Preconditions

- `github-fetch` binary is built and on `req.Binary`.

## Steps

1. Set `req.Args = ["install"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"install"}
	return nil
}
```
