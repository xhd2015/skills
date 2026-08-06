# Scenario

**Feature**: multi-skill Registry lists, shows, installs, and batch-updates

```
# registry host
caller -> Registry.HandleSkill / HandleSkills -> list | show | install | update
```

## Preconditions

- RegistrySkills is an ordered slice of RegisteredSkill entries.

## Steps

1. Set `req.Mode = ModeRegistry` with foo/bar skills.
2. Leaves choose RegistryCmd and Args; update leaves pre-install.

## Context

- Display names for update messages use RegisteredSkill.Name.

```go
import "github.com/xhd2015/skills/skillcmd"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Mode = ModeRegistry
	req.RegistryCmd = RegistryCmdSkill
	req.RegistrySkills = []skillcmd.RegisteredSkill{
		{Name: "foo", Description: "foo skill description", Content: fooSkillContent},
		{Name: "bar", Description: "bar skill description", Content: barSkillContent},
	}
	return nil
}
```
