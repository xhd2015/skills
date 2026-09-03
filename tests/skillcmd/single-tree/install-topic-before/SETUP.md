# Scenario

**Feature**: show-style `skill <topic> --install` peels the topic; install uses destination

```
caller -> SingleSkill.Handle(skill-cli --install --dir vendor/skills --dry-run)
  -> peels topic skill-cli
  -> installs whole skill under vendor/skills/demo-skill
```

## Preconditions

- Parent single-tree TreeFS includes `skill-cli` topic.
- UseWorkDir for filesystem isolation.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.UseWorkDir = true
	return nil
}
```
