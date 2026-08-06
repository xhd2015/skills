# Scenario

**Feature**: skillcmd.HandleInstall performs default-target installs

```
# install API lives in skillcmd (install package becomes a shim)
caller -> skillcmd.HandleInstall(opts, args) -> .agents/skills/<name>/
```

## Preconditions

- Isolated work directory for filesystem side effects.

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
