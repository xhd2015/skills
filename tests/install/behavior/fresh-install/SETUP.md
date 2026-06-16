## Preconditions
- The target directory "example-skill" does **not** exist in the working directory.

## Steps
1. Call `HandleInstall` with args `["example-skill"]` and fresh skill content `"# test skill\n"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"example-skill"}
	return nil
}
```
