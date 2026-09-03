# Scenario

**Feature**: skillcmd.HandleInstall default and explicit `--dir` / `<dir>` targets

```
# install API lives in skillcmd (install package becomes a shim)
caller -> skillcmd.HandleInstall(opts, args) -> .agents/skills/<name>/
caller -> HandleInstall(--dir <collection>) -> <collection>/<name>/
```

## Preconditions

- Isolated work directory for filesystem side effects.
- Explicit `--dir` / positional `<dir>` use smart layout (basename `skills`
  nests; existing `SKILL.md` or matching basename is direct).

## Steps

1. Set ModeInstallCompat and UseWorkDir.
2. Leaf configures InstallOpts and Args.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Mode = ModeInstallCompat
	req.UseWorkDir = true
	return nil
}
```
