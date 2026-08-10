# Scenario

**Feature**: eval flag argument validation

```
user -> playwright-debug CLI (-e|--eval bad usage) -> routing error
```

## Steps

1. Each leaf sets invalid eval flag combinations.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```