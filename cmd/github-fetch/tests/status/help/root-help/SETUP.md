# Scenario

**Feature**: root help lists the status command

```
# root command help
github-fetch -h -> command list includes status
```

## Preconditions

- Root help text is updated to document `status`.

## Steps

1. Set `req.Args = []string{"-h"}`.
2. Run and assert `status` appears in command list.

## Context

- Does not run the status logic; only checks root help wiring.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"-h"}
	return nil
}
```