# Scenario

**Feature**: update reports `skill not installed` for missing registry skills

```
# no pre-install -> both skills missing
caller -> HandleSkills(update) -> foo/bar  not installed + summary
```

## Steps

1. Do not pre-install any skill.
2. Run `skills update` (Args set by parent).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.PreInstall = false
	return nil
}
```
