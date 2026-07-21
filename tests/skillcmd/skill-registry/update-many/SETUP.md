# Scenario

**Feature**: batch update refreshes installed skills and reports missing ones

```
# skills update walks registry
caller -> Registry.HandleSkills(update ...) -> up-to-date | skill not installed
```

## Preconditions

- Uses RegistryCmdSkills with update subcommand.
- Workdir isolation required for install probes.

## Steps

1. Set RegistryCmdSkills, UseWorkDir, Args for update.
2. Leaves pre-install and/or mutate as needed.

```go
func Setup(t *testing.T, req *Request) error {
	req.RegistryCmd = RegistryCmdSkills
	req.UseWorkDir = true
	req.Args = []string{"update"}
	return nil
}
```
