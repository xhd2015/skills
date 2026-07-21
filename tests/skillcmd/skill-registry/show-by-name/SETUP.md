# Scenario

**Feature**: show accepts skill name before or after `--show`

```
# both argument orders yield the same skill body
caller -> HandleSkill(--show foo | foo --show) -> foo content
```

## Preconditions

- Registry includes skill `foo`.

## Steps

1. Leaves set Args for each flag order.

```go
func Setup(t *testing.T, req *Request) error {
	req.RegistryCmd = RegistryCmdSkill
	return nil
}
```
