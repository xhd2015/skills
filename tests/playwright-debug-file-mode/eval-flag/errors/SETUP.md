# Scenario

**Feature**: eval flag argument validation

```
user -> playwright-debug CLI (-e|--eval bad usage) -> routing error
```

## Steps

1. Each leaf sets invalid eval flag combinations.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```