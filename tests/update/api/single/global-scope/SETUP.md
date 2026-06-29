# Scenario

**Feature**: `--global` resolves targets under `$HOME`

```
HandleInstall --global / HandleUpdate --global -> home-relative skill paths
```

## Preconditions

- Leaves set `req.UseGlobalHome = true`.

## Steps

1. Use `--global` on install and update flag lists.

## Context

- Project-local default paths must remain absent when only global install ran.

```go
func Setup(t *testing.T, req *Request) error {
	req.UseGlobalHome = true
	return nil
}
```