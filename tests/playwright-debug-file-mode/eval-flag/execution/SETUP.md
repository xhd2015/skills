# Scenario

**Feature**: `-e` runs adhoc eval snippet

```
user -> playwright-debug CLI (-e '<script>') -> eval runner -> eval-ok
```

## Context

- Labeled `slow` because eval still launches playwright.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```