# Scenario

**Feature**: `-e` without script argument fails

```
user -> playwright-debug CLI (-e) -> requires script error
```

## Steps

1. Set `req.Args = []string{"-e"}`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"-e"}
	return nil
}
```