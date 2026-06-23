# Scenario

**Feature**: `skill show` output scope depends on the `--header` flag

```
# header flag selects frontmatter-only output
user -> skill CLI (skill show --header) -> delimited YAML block

# default show prints full SKILL.md
user -> skill CLI (skill show) -> header + body
```

## Preconditions

- Root setup built the `go-best-practice` binary on `req.Binary`.

## Steps

1. Each leaf sets `req.HeaderOnly` for its branch.
2. Run invokes `skill show` with or without `--header`.

## Context

- Body marker `# Go Best Practice Skill` distinguishes body content from header metadata.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Binary == "" {
		t.Fatal("req.Binary must be set by root setup")
	}
	return nil
}
```