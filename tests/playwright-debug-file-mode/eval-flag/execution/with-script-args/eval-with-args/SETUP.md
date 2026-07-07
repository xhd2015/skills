# Scenario

**Feature**: `-e` forwards trailing script arguments

```
user -> playwright-debug CLI (-e '<script>' baz) -> ["baz"]
```

## Steps

1. Set `req.Args` with `-e`, argv-printing script, and `baz`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-e", `console.log(JSON.stringify(process.argv.slice(3)))`, "baz"}
	return nil
}
```