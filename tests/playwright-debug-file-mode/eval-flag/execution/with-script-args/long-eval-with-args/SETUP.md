# Scenario

**Feature**: `--eval` forwards multiple trailing script arguments

```
user -> playwright-debug CLI (--eval '<script>' a b) -> ["a","b"]
```

## Steps

1. Set `req.Args` with `--eval`, argv-printing script, and `a`, `b`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--eval", `console.log(JSON.stringify(process.argv.slice(3)))`, "a", "b"}
	return nil
}
```