# Scenario

**Feature**: planned set or file content differs from disk

```
# content update, dropped nested path, or rename SKILL.md -> TOPIC.md
HandleInstall -> Update skill at + create:/update:/delete: lines
# incremental: only touched paths; no whole-dir RemoveAll as the normal path
```

## Preconditions

- Skill directory already exists.
- Plan differs from disk by content and/or membership.

## Steps

1. Leaves pre-seed disk state and set plan (SkillContent + ExtraFiles).
2. Apply install; assert per-file stdout actions and resulting files.

## Context

- Rename is reported as `delete:` old path + `create:` new path (no `rename:`).
- Detail order: creates (sorted), updates (sorted), deletes (sorted).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Plan-changed leaves always start from an existing skill directory.
	req.PreExistingDir = "example-skill"
	return nil
}
```
