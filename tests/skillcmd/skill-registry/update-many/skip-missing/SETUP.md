# Scenario

**Feature**: update reports `skill not installed` for missing registry skills

```
# no pre-install -> both skills missing
caller -> HandleSkills(update) -> skill not installed: foo / bar
```

## Steps

1. Do not pre-install any skill.
2. Run `skills update` (Args set by parent).

```go
func Setup(t *testing.T, req *Request) error {
	req.PreInstall = false
	return nil
}
```
