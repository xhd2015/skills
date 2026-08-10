# Scenario

**Feature**: disk matches plan with no unplanned files

```
# pre-seed skill dir with exact plan files
pre-seed example-skill/{SKILL.md,...} == plan
HandleInstall -> "Skill is up to date" (no writes)
```

## Preconditions

- On-disk skill dir already contains exactly the planned regular files with
  matching content (no orphans).

## Steps

1. Leaves pre-seed PreExistingDir + files matching plan.
2. Call HandleInstall (apply or `--dry-run`).

## Context

- Up-to-date path must not print create/update/delete detail lines.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Clean-match leaves seed disk equal to plan; default target is explicit dir.
	req.SkillContent = "# test skill\n"
	return nil
}
```
