# Scenario

**Feature**: `--global` batch update with zero installs reports each skill

```
HOME temp dir, no installs -> HandleUpdateMany(..., --global) -> not-installed line per skill
```

## Preconditions

- `req.UseGlobalHome` isolates `HOME` with no skill trees.

## Steps

1. Enable global home temp dir.
2. Run batch update with `["--global"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.UseGlobalHome = true
	req.PreInstalls = nil
	req.UpdateArgs = []string{"--global"}
	return nil
}
```