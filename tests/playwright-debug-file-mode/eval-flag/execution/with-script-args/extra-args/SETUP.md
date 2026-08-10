# Scenario

**Feature**: `--eval` forwards trailing script argument (formerly rejected)

```
user -> playwright-debug CLI (--eval '<script>' extra) -> ["extra"]
```

## Steps

1. Set `req.Args` with `--eval`, argv-printing script, and `extra`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--eval", `console.log(JSON.stringify(process.argv.slice(3)))`, "extra"}
	return nil
}
```