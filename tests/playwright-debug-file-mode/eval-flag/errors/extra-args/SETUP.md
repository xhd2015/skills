# Scenario

**Feature**: `--eval` rejects extra positional arguments

```
user -> playwright-debug CLI (--eval 'x' extra) -> unexpected arguments error
```

## Steps

1. Set `req.Args = []string{"--eval", "x", "extra"}`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--eval", "x", "extra"}
	return nil
}
```