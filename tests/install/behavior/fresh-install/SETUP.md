# Scenario

**Feature**: fresh install when target directory is missing

```
# no example-skill/ on disk
HandleInstall(example-skill) -> "Installed skill to: <absDir>"
```

## Preconditions
- The target directory "example-skill" does **not** exist in the working directory.

## Steps
1. Call `HandleInstall` with args `["example-skill"]` and fresh skill content `"# test skill\n"`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"example-skill"}
	return nil
}
```
