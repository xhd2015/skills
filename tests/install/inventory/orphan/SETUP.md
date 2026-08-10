# Scenario

**Bug**: unplanned on-disk files must prevent up-to-date and be deleted

```
# plan matches root SKILL.md but orphan.txt remains on disk
example-skill/{SKILL.md match, orphan.txt}
HandleInstall -> not up to date; delete: <abs>/orphan.txt
```

## Preconditions

- Planned files match content; at least one regular file under the skill dir is
  not in the plan (orphan).

## Steps

1. Leaves seed matching plan files plus orphan path(s).
2. Apply or dry-run install with the same plan.

## Context

- Orphans include leftover nested `SKILL.md` after renames to `TOPIC.md`.
- Pure cleanup still uses header `Update skill at` (dir already existed).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Orphan leaves keep planned root content matching disk; only unplanned files differ.
	req.SkillContent = "# test skill\n"
	return nil
}
```
