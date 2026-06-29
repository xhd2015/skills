# Scenario

**Feature**: choose single-skill or registry batch update API

```
# single skill definition
test -> HandleUpdate(opts, args) -> per-target InstallTo

# many registry entries share one flag parse
test -> HandleUpdateMany(skills, args) -> per-skill, per-target InstallTo
```

## Preconditions

- Descendants set `req.UseMany` and populate either `SingleOpts` or `ManySkills`.

## Steps

1. Descendant grouping nodes set the API mode on `Request`.

## Context

- API choice is the highest-impact split in this tree.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	return nil
}
```