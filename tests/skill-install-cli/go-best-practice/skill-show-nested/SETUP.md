# Scenario

**Feature**: nested topic show via skill --show path (Shape 3)

```
# both flag orders for nested path
user -> go-best-practice skill --show skill-cli -> nested skill-cli body
user -> go-best-practice skill skill-cli --show -> same content
```

## Preconditions

- go-best-practice embeds nested skill-cli/SKILL.md after migration.

## Steps

1. Leaves set Args for each flag order.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Binary == "" {
		t.Fatal("req.Binary must be set")
	}
	return nil
}
```
