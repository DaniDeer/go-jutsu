# Golden Files Pattern

Store expected output in files. Compare actual vs expected. Update with `-update` flag.

## What It Is (and Isn't)

Snapshot testing. Test fixtures in files. Easy to review output changes.

Not for simple values. Not database fixtures. Pattern for complex output (HTML, JSON, etc).

## Where You See It

**Template rendering:**

```go
func TestRenderTemplate(t *testing.T) {
    got := RenderTemplate(data)
    golden := filepath.Join("testdata", "output.golden")

    if *update {
        os.WriteFile(golden, []byte(got), 0644)
    }

    want, _ := os.ReadFile(golden)
    if got != string(want) {
        t.Errorf("mismatch")
    }
}
```

**Compiler output:**

```go
// Test AST, bytecode, error messages
```

## Real Example

````go
var update = flag.Bool("update", false, "update golden files")

func TestFormatJSON(t *testing.T) {
    tests := []struct{
        name string
        input map[string]interface{}
    }{
        {"simple", map[string]interface{}{"key": "value"}},
        {"nested", map[string]interface{}{"a": map[string]int{"b": 1}}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FormatJSON(tt.input)
            goldenFile := filepath.Join("testdata", tt.name+".golden")
            
            if *update {
                os.WriteFile(goldenFile, []byte(got), 0644)
                return
            }
            
            want, err := os.ReadFile(goldenFile)
            if err != nil {
                t.Fatalf("failed to read golden file: %v", err)
            }
            
            if got != string(want) {
                t.Errorf("output mismatch\
got:\
%s\
want:\
%s", got, want)
            }
        })
    }
}
```

## Gotchas

**Update golden files:**
```sh
go test -update
```

**testdata directory:**
```
package/
  code.go
  code_test.go
  testdata/
    output.golden
    error.golden
```

**Version control:**
- Commit golden files
- Review changes in PR
- Shows exact output differences

````
