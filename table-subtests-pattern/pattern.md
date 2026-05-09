# Table Subtests Pattern

Run same test with different inputs. Use `t.Run()` for isolation. Parallel execution support.

## What It Is (and Isn't)

Data-driven testing. Each case is subtest. Clear pass/fail per scenario.

Not copy-paste tests. Not single assertion. Structured test organization.

## Where You See It

**Input/output pairs:**

```go
tests := []struct{
    name string
    input int
    want int
}{
    {"zero", 0, 0},
    {"positive", 5, 25},
    {"negative", -3, 9},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := Square(tt.input)
        if got != tt.want {
            t.Errorf("got %d want %d", got, tt.want)
        }
    })
}
```

## Real Example

```go
func TestParseURL(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
        wantHost string
    }{
        {
            name: "valid http",
            input: "http://example.com",
            wantErr: false,
            wantHost: "example.com",
        },
        {
            name: "invalid scheme",
            input: "ftp://example.com",
            wantErr: true,
        },
        {
            name: "missing scheme",
            input: "example.com",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()  // Run in parallel

            got, err := ParseURL(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && got.Host != tt.wantHost {
                t.Errorf("host = %v, want %v", got.Host, tt.wantHost)
            }
        })
    }
}
```

## Gotchas

**Parallel tests need t.Parallel():**

```go
t.Run(name, func(t *testing.T) {
    t.Parallel()  // Each subtest runs concurrently
    // ...
})
```

**Loop variable capture:**

```go
for _, tt := range tests {
    tt := tt  // Go 1.21 and earlier
    t.Run(tt.name, func(t *testing.T) {
```

**Run specific subtest:**

```sh
go test -run TestName/subtest_name
```
