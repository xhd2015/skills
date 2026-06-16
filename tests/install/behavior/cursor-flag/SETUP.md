## Preconditions
- No pre-existing target directory.
- The `--cursor` flag is passed (without `--global`).

## Steps
1. Call `HandleInstall` with `--cursor` to install to `.cursor/skills/<skillDirName>`.

## Context
- The `--cursor` flag installs to `.cursor/skills/<skillDirName>` in the current working directory.
- This test verifies it works without `--global` (local install).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--cursor"}
	return nil
}
```
