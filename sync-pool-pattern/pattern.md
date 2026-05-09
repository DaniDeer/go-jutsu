# sync.Pool Pattern

Reuse objects to reduce GC pressure. Temporary object cache automatically cleaned by GC.

## What It Is (and Isn't)

Object pool managed by runtime. Reduces allocations. No size limit or lifecycle control.

Not long-term storage. Not guaranteed retention. Pool may be cleared anytime.

## Where You See It

**Buffer reuse:**

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

buf := bufferPool.Get().(*bytes.Buffer)
buf.Reset()
defer bufferPool.Put(buf)
```

**fmt package:**

```go
// fmt uses sync.Pool for printers internally
fmt.Printf("%s", value)
```

## Real Example

```go
var requestPool = sync.Pool{
    New: func() interface{} {
        return &Request{}
    },
}

func HandleRequest(data []byte) {
    req := requestPool.Get().(*Request)
    req.Reset()
    defer requestPool.Put(req)

    req.Parse(data)
    req.Process()
}
```

## Gotchas

**Reset objects:**

```go
buf := pool.Get().(*bytes.Buffer)
buf.Reset()  // MUST reset before use
```

**Type assertion:**

```go
obj := pool.Get().(*MyType)  // Panic if wrong type
```

**No guarantees:**

- Object may not be in pool (New called)
- Pool cleared between GC runs
- Use for short-lived objects only
