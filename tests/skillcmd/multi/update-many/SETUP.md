# Scenario

**Feature**: batch update refreshes installed skills and reports missing ones

```
# skills update walks registry
caller -> Registry.HandleSkills(update ...) -> polished status | not installed + summary
```

## Preconditions

- Uses RegistryCmdSkills with update subcommand.
- Workdir isolation required for install probes.

## Steps

1. Set RegistryCmdSkills, UseWorkDir, Args for update.
2. Leaves pre-install and/or mutate as needed.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.RegistryCmd = RegistryCmdSkills
	req.UseWorkDir = true
	req.Args = []string{"update"}
	return nil
}
```
