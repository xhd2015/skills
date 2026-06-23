## Expected

- `GetHeader` succeeds.
- `ParseHeader` succeeds.
- `Entries.Get("description")` returns `multi line` as a single normalized value.

## Side Effects

- None.

## Errors

- No errors from `GetHeader` or `ParseHeader`.

## Exit Code

- Not applicable.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.HeaderErr != nil {
		t.Fatalf("GetHeader error: %v", resp.HeaderErr)
	}
	if resp.ParseErr != nil {
		t.Fatalf("ParseHeader error: %v", resp.ParseErr)
	}
	if got := resp.Entries.Get("description"); got != "multi line" {
		t.Fatalf(`Entries.Get("description") = %q, want "multi line"`, got)
	}
}
```